package controllers

import (
	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/revel/revel"
)

var (
	userPermission = []string{"Problems.New"}
)

func inStringSlice(s string, slc []string) bool {
	for _, v := range slc {
		if s == v {
			return true
		}
	}
	return false
}

type App struct {
	*revel.Controller
}

func connected(c *revel.Controller) bool {
	username, has := c.Session["username"]
	if !has {
		return false
	}
	u := new(models.User)
	if has, err := engine.Where("name = ?", username).Get(u); has && err == nil {
		return true
	} else {
		println(has)
		println(err)
		return false
	}
}
func (c App) Index() revel.Result {

	return c.Render()
}

func CheckLogin(c *revel.Controller) revel.Result {
	println(c.Action)
	if inStringSlice(c.Action, userPermission) {
		if ok := connected(c); !ok {
			c.Flash.Error("请先登陆")
			return c.Redirect(routes.Problems.Index())
		} else {
			return nil
		}
	}
	return nil
}
