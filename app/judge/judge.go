package judge

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
)

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

//judge input and output
func Judge(language string, filePath, inputPath, outputPath string, timeLimit, memoryLimit int64) (int, error) {
	defer os.Remove(filePath + "/tmp")
	cmd := exec.Command("sandbox", "--lang="+language, "--time="+strconv.FormatInt(timeLimit, 10), "--memory="+strconv.FormatInt(memoryLimit, 10), filePath+"/tmp."+language, filePath+"/tmp", inputPath, outputPath)
	fmt.Println(cmd.Args)
	testOut, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", testOut)
	if err != nil {
		return models.WrongAnswer, err
	}
	if fmt.Sprintf("%s", testOut) == "AC" {
		return models.Accept, nil
	}
	if fmt.Sprintf("%s", testOut) == "TLE" {
		return models.TimeLimitExceeded, nil
	}
	if fmt.Sprintf("%s", testOut) == "MLE" {
		return models.MemoryLimitExceeded, nil
	}
	if fmt.Sprintf("%s", testOut) == "CE" {
		return models.CompileError, nil
	}
	return models.WrongAnswer, nil
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

func HandleCodeLoop() {
	for sources := range unHandledCodeChan {
		fmt.Println("update")
		for _, v := range sources {
			problem := new(models.Problem)
			_, err := engine.Id(v.ProblemId).Cols("input_test_path", "output_test_path", "time_limit", "memory_limit").Get(problem)
			if err != nil {
				fmt.Println(err)
			}
			result, err := Judge(v.LangString(), v.Path, problem.InputTestPath, problem.OutputTestPath, problem.TimeLimit, problem.MemoryLimit)
			if err != nil {
				panic(err)
			} else {
				v.Status = result
				if result == models.Accept {
					n, err := engine.Id(v.Id).Cols("status").Update(&v)
					if err != nil {
						fmt.Println(n)
						panic(err)
					}
					p := new(models.Problem)
					_, err = engine.Id(v.ProblemId).Incr("solved", 1).Update(p)
					if err != nil {
						fmt.Println(err)
					}
					u := new(models.User)
					_, err = engine.Id(v.UserId).Incr("solved", 1).Update(u)
				} else {
					n, err := engine.Id(v.Id).Cols("status").Update(&v)
					if err != nil {
						fmt.Println(n)
						panic(err)
					}
				}
			}
		}
	}
}
