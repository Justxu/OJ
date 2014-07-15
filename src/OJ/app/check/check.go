package check

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"OJ/app/models"

	"github.com/ggaaooppeenngg/util"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

const (
	imageName = "ubuntu/sandbox"
)

func init() {
	engine = models.Engine()
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

//check input and output
func CheckInput(language string, filePath, inputPath, outputPath string) (int, error) {
	inf, err := os.Open(inputPath)
	if err != nil {
		return models.UnHandled, err
	}
	cmd := exec.Command("sandbox", "--glang="+language, filePath+"/tmp."+language, filePath+"tmp")
	cmd.Stdin = inf
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return models.UnHandled, err
	}
	outf, err := os.Open(outputPath)
	if err != nil {
		return models.UnHandled, err
	}
	var testOut []byte
	tmp := make([]byte, 256)
	for n, err := outf.Read(tmp); err != io.EOF; n, err = outf.Read(tmp) {
		testOut = append(testOut, tmp[:n]...)
	}
	if bytes.Equal(out.Bytes(), testOut) {
		return models.Accept, nil
	} else {
		return models.WrongAnswer, nil
	}
}

//生产者要不断地扫描任务
//但是在任务处理完成之前，还是保持着未处理状态，会被再次加入任务队列这个时候
//有一个通道用来取任务
//每个处理线程都要有个
//一种方法每次扫描完成后的所有任务完成再通知才让生产者继续扫描
//所以需要两个通道，一个用于缓冲任务，一个用于告知结束,select或许可以用上来

func Do() {
	var sources []models.Source
	var problem models.Problem
	err := engine.Where("status = ?", models.UnHandled).Find(&sources)
	fmt.Println(len(sources))
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range sources {
		engine.Id(v.ProblemId).Cols("input_test", "output_test").Get(&problem)
		result, err := CheckInput(v.Lang, v.Path, problem.InputTest, problem.OutputTest)
		if err != nil {
			panic(err)
		} else {
			if result == models.Accept {
				v.Status = models.Accept
				engine.Id(v.Id).Cols("status").Update(&v)
			} else {
				v.Status = models.WrongAnswer
				engine.Id(v.Id).Cols("status").Update(&v)
				engine.Id(v.ProblemId).Incr("solved = solved + ?", 1)
			}
		}
	}
	fmt.Println("refresh")
}
