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

func Do() {
	var sources []models.Source
	err := engine.Where("status = ?", models.UnHandled).Find(&sources)
	fmt.Println(len(sources))
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range sources {
		//util.Run("docker", "build", "-t", imageName, ".")
		//util.Run("docker", imageName, "go", "run", "main.go")
		out, err := util.Run("go", "run", v.Path)
		if err != nil {
			fmt.Println(err)
		}
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
