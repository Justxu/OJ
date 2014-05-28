package check

import (
	"fmt"

	"OJ/app/models"

	"github.com/ggaaooppeenngg/util"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

const (
	imageName = "ggaaooppeenngg/ubuntu:golang"
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
}

//clean the container after running
func removeContainer() {
	util.Run("docker", "rm", "check")
}

//user generated dockfile to build and run the test
func test(path string) {
	defer removeContainer()
	genDocFile(path)
	_, err := util.Run("docker", "build", "-t", imageName, ".")
	if err != nil {
		fmt.Println(err)
	}
	out, err := util.Run("docker", "run", "--name=check", imageName, "go", "run", "/home/main.go")
	if err != nil {
		fmt.Println(err)
	}
}

func Do() {
	var sources []models.Source
	err := engine.Where("status = ?", models.UnHandled).Find(&sources)
	fmt.Println(len(sources))
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range sources {
		out, _ := util.Run("go", "run", v.Path)
		if string(out) != "Hello World\n" {
			v.Status = models.WrongAnswer
			engine.Id(v.Id).Cols("status").Update(&v)
		} else {
			v.Status = models.Accept
			engine.Id(v.Id).Cols("status").Update(&v)
		}
	}
	fmt.Println("refresh")
}
