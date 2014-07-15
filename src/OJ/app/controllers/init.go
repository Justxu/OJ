package controllers

import (
	"fmt"
	//	"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/revel/config"
	//	"github.com/revel/revel/modules/jobs/app/jobs"

	//	"OJ/app/check"
	"OJ/app/models"
)

var (
	engine *xorm.Engine
	stmp   SmtpConfig
)

func init() {
	fmt.Println("init")
	engine = models.Engine()
	/*	revel.OnAppStart(func() {
			jobs.Every(time.Second, jobs.Func(check.Do))
		})
	*/
	c, err := config.ReadDefault("src/OJ/conf/misc.conf")
	if err != nil {
		panic(err)
	}
	stmp.Username = c.String("smtp", "username")
	stmp.Password = c.String("smtp", "password")
	stmp.Host = c.String("smtp", "host")
	stmp.Addr = c.String("smtp", "address")
}
