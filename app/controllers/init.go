package controllers

import (
	//"fmt"
	//"time"

	"github.com/go-xorm/xorm"
	"github.com/revel/config"
	"github.com/revel/revel"

	//"OJ/app/check"
	"github.com/ggaaooppeenngg/OJ/app/models"
)

var (
	admin      string
	engine     *xorm.Engine
	smtpConfig SmtpConfig
)

func GetStmp() SmtpConfig {
	return smtpConfig
}

func initTemplateFunc() {
	revel.TemplateFuncs["isAdmin"] = func(a interface{}) bool {
		if a == nil {
			return false
		} else {
			println(admin)
			return a.(string) == admin
		}
	}
}

//check permission
func initIntercepter() {
	revel.InterceptFunc(CheckLogin, revel.BEFORE, &Problems{})
}

func init() {
	engine = models.Engine()
	/*	revel.OnAppStart(func() {
			jobs.Every(time.Second, jobs.Func(check.Do))
		})
	*/
	c, err := config.ReadDefault("conf/misc.conf")
	if err != nil {
		panic(err)
	}
	smtpConfig.Username, err = c.String("smtp", "username")
	if err != nil {
		panic(err)
	}
	smtpConfig.Password, _ = c.String("smtp", "password")
	smtpConfig.Host, _ = c.String("smtp", "host")
	smtpConfig.Addr, _ = c.String("smtp", "address")
	admin, err = c.String("app", "admin")
	if err != nil {
		panic(err)
	}
	initIntercepter()
	initTemplateFunc()

}
