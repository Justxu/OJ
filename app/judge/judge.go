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
	unHandledCodeChan = make(chan []models.Source, 32)
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
	testOut, err := cmd.CombinedOutput()
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

//生产者要不断地扫描任务
//但是在任务处理完成之前，还是保持着未处理状态，会被再次加入任务队列这个时候
//有一个通道用来取任务
//每个处理线程都要有个
//一种方法每次扫描完成后的所有任务完成再通知才让生产者继续扫描
//所以需要两个通道，一个用于缓冲任务，一个用于告知结束,select或许可以用上来

func GetHandledCodeLoop() {
	var sources []models.Source
	for {
		time.Sleep(time.Second)
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
		fmt.Println("refresh")
		unHandledCodeChan <- sources
	}
}

func HandleCodeLoop() {
	for sources := range unHandledCodeChan {
		for _, v := range sources {
			problem := new(models.Problem)
			_, err := engine.Id(v.ProblemId).Cols("input_test_path", "output_test_path").Get(problem)
			if err != nil {
				fmt.Println(err)
			}
			result, err := Judge(v.LangString(), v.Path, problem.InputTestPath, problem.OutputTestPath, problem.TimeLimit, problem.MemoryLimit)
			if err != nil {
				panic(err)
			} else {
				v.Status = result
				if result == models.Accept {
					engine.Id(v.Id).Cols("status").Update(&v)
					engine.Id(v.ProblemId).Incr("solved = solved + ?", 1)
				} else {
					engine.Id(v.Id).Cols("status").Update(&v)
				}
			}
			fmt.Println("update")
		}
	}
}
