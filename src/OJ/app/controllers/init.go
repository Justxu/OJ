package controllers

import (
	"fmt"
	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel"
	"github.com/revel/revel/modules/jobs/app/jobs"

	"OJ/app/check"
	"OJ/app/models"
)

var (
	engine *xorm.Engine
)

func init() {
	fmt.Println("init")
	engine = models.Engine()
	revel.OnAppStart(func() {
		jobs.Every(time.Second, jobs.Func(check.Do))
	})
}
