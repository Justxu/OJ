package check

import (
	"fmt"
	"time"

	"OJ/app/models"

	"github.com/ggaaooppeenngg/util"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func init() {
	engine = models.Engine()
}
func Do() {
	for {
		var sources []models.Source
		engine.Where("status = ?", models.UnHandled).Find(&sources)
		for _, v := range sources {
			out, _ := util.Run("go", "run", v.Path)
			if string(out) != "Hello World\n" {
				v.Status = models.WrongAnswer
				engine.Update(v)
			} else {
				v.Status = models.Accept
				engine.Update(v)
			}
		}
		fmt.Println("refresh")
		time.Sleep(time.Second)
	}
}
