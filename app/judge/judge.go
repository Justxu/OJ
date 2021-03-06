package judge

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ggaaooppeenngg/OJ/app/models"

	"github.com/ggaaooppeenngg/util"
	"github.com/go-xorm/xorm"
)

var (
	engine *xorm.Engine

	unHandledCodeChan chan []models.Source
)

const (
	imageName = "ubuntu/sandbox"

	seperator = "!-_-"
)

type result struct {
	Status      int
	Time        int64
	Memory      int64
	Nth         int
	WrongAnswer string
	PanicOutput string //异常的输出
}

func init() {
	engine = models.Engine()
	//buffer of size 32
	unHandledCodeChan = make(chan []models.Source)
}

//use command sed repleace SRCFILE to real source file
func genDocFile(path string) error {
	out, err := util.Run("sed", "s/SRCFILE/"+path+"/g", "Seedfile")
	if err != nil {
		return err
	}
	_, err = util.WriteFile("Dockerfile", out)
	if err != nil {
		return err
	}
	return nil
}

//clean the container after running
func removeContainer(name string) {
	util.Run("docker", "rm", name)
}

//user generated dockfile to build and run the test
func test(path string) []byte {
	defer removeContainer(path)
	genDocFile(path)
	_, err := util.Run("docker", "build", "-t", path, ".")
	if err != nil {
		fmt.Println(err)
	}
	out, err := util.Run("docker", "run", "-i", "--name="+path, imageName, "go", "run", "/home/main.go")
	if err != nil {
		fmt.Println(err)
	}
	return out
}

/*
 Every single test is seperated with a new line "!-_-"
 like a+b :
 input:		output:
 2
 1 2		3
 3 4		7
 !-_-		!-_-
 2
 2 1		3
 4 3		7
 !-_-		!-_-
 2
 3 4		7
 1 2		3
*/

//judge input and output
func Judge(language string, filePath, inputPath, outputPath string, timeLimit, memoryLimit int64) *result {
	defer os.Remove(filePath + "/tmp")
	cmd := exec.Command("sandbox", "--lang="+language, "--time="+strconv.FormatInt(timeLimit, 10), "--memory="+strconv.FormatInt(memoryLimit, 10), "-c", "-s", filePath+"/tmp."+language, "-b", filePath+"/tmp", "-i", inputPath, "-o", outputPath)
	testOut, err := cmd.CombinedOutput()
	fmt.Println(cmd.Args)
	fmt.Printf("%s\n", testOut)
	if err != nil {
		fmt.Println(err)
		return &result{Status: models.PanicError, PanicOutput: err.Error()}
	}
	//panic output
	if !regexp.MustCompile(`\w\w:\d+:\d+:([\s\S]*)`).Match(testOut) {
		fmt.Printf("out put format error:\n%s\n", testOut)
		return &result{Status: models.PanicError, PanicOutput: fmt.Sprintf("%s", testOut)}
	}
	return getResults(testOut)
}

/*
	use producer-consumer pattern to handle codes.
	pick up "unhandled" codes and immediately update to "handling" status,
	and send them to "unHandledCodeChan" for consumer to deal with

	//TODO: use consumber goroutine pool to handle codes in order to increase multithreading degreee
*/
func GetHandledCodeLoop() {
	for {
		fmt.Println("refresh")
		time.Sleep(2 * time.Second)
		sources := make([]models.Source, 0)
		err := engine.Where("status = ?", models.UnHandled).Find(&sources)
		if len(sources) == 0 {
			continue
		}
		if err != nil {
			fmt.Println(err)
		}
		//TODO:if panics rollback the status
		_, err = engine.Where("status = ?", models.UnHandled).Cols("status").Update(&models.Source{Status: models.Handling})
		if err != nil {
			fmt.Println(err)
		}
		if len(sources) > 0 {
			fmt.Println("produce")
			fmt.Println(sources[0].StatusString())
			unHandledCodeChan <- sources
		}
	}
}

func getResults(out []byte) *result {
	results := strings.Split(fmt.Sprintf("%s", out), ":")
	statuss := results[0]
	var status int
	var wrongAnswer string
	var memory, time, nth64 int64
	if len(results) > 1 {
		memory, _ = strconv.ParseInt(results[1], 0, 64)
	}
	if len(results) > 2 {
		time, _ = strconv.ParseInt(results[2], 0, 64)
	}
	if len(results) > 3 {
		nth64, _ = strconv.ParseInt(results[3], 0, 64)
	}
	switch statuss {
	case "AC":
		status = models.Accept
	case "CE":
		status = models.CompileError
	case "TL":
		status = models.TimeLimitExceeded
	case "ML":
		status = models.MemoryLimitExceeded
	case "RE":
		status = models.RuntimeError
	case "FE":
		status = models.PresentationError
	case "WA":
		status = models.WrongAnswer
		if len(results) > 4 {
			wrongAnswer = results[4]
		}
		nth64 += 1
	}
	return &result{Status: status, Time: time, Memory: memory, Nth: int(nth64), WrongAnswer: wrongAnswer}
}

//add rollbcak when panic
func HandleCodeLoop() {
	for sources := range unHandledCodeChan {
		fmt.Println("update")
		for _, v := range sources {
			problem := new(models.Problem)
			_, err := engine.Id(v.ProblemId).Cols("input_test_path", "output_test_path", "time_limit", "memory_limit").Get(problem)
			if err != nil {
				fmt.Println(err)
			}
			result := Judge(v.LangString(), v.Path, problem.InputTestPath, problem.OutputTestPath, problem.TimeLimit, problem.MemoryLimit)
			if result.Status == models.PanicError {
				//if panic err,it is caused by sandbox errors
				v.Status = models.PanicError
				v.PanicError = result.PanicOutput
				_, err := engine.Id(v.Id).Cols("status,panic_error").Update(&v)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				v.Status = result.Status
				v.Time = time.Duration(result.Time) / time.Millisecond
				v.Memory = result.Memory
				v.Nth = result.Nth
				v.WrongAnswer = result.WrongAnswer
				transaction := engine.NewSession()
				if v.Status == models.Accept {
					_, err := transaction.Id(v.Id).Cols("status", "time", "memory", "nth").Update(&v)
					if err != nil {
						//if panic error,it is caused by database design,like interati constraints
						err = transaction.Rollback()
						if err != nil {
							fmt.Println(err)
						}
						v.Status = models.PanicError
						_, err = engine.Id(v.Id).Cols("status").Update(&v)
						if err != nil {
							fmt.Println(err)
						}
					} else {
						p := new(models.Problem)
						_, err = engine.Id(v.ProblemId).Incr("solved", 1).Update(p)
						if err != nil {
							fmt.Println(err)
						}
						u := new(models.User)
						_, err = engine.Id(v.UserId).Incr("solved", 1).Update(u)
					}
				} else {
					_, err := engine.Id(v.Id).Cols("status", "time", "memory", "nth", "wrong_answer").Update(&v)
					//if panic error,it is caused by database design,like interati constraints
					if err != nil {
						err = transaction.Rollback()
						if err != nil {
							fmt.Println(err)
						}
						v.Status = models.PanicError
						_, err = engine.Id(v.Id).Cols("status").Update(&v)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}
}
