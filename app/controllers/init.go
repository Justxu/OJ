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
	appAddr    string
	appPort    string
	engine     *xorm.Engine
	smtpConfig SmtpConfig
)

func GetStmp() SmtpConfig {
	return smtpConfig
}

func initTemplateFunc() {
	revel.TemplateFuncs["isAdmin"] = IsAdmin
	revel.TemplateFuncs["Text"] = Text
}

func initIntercepter() {
	revel.InterceptMethod(CheckLogin, revel.BEFORE)
}

func init() {
	engine = models.Engine()
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
	c, err = config.ReadDefault("conf/app.conf")
	if err != nil {
		panic(err)
	}
	appAddr, err = c.String("", "http.addr")
	appPort, err = c.String("", "http.port")
	if appAddr == "" || appPort == "" {
		panic("init fail,can not get address and port from conf/misc.conf")
	}
	initIntercepter()
	initTemplateFunc()

}
