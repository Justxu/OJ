package controllers

import (
	"fmt"
	"time"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"

	"OJ/app/models"
	"OJ/app/routes"
)

type Code struct {
	*revel.Controller
}

func (c *Code) Answer(id int64) revel.Result {
	var problem models.Problem
	engine.Id(id).Get(&problem)
	return c.Render(problem)
}

func (c *Code) Submit(code string, problemId int64, lang string) revel.Result {
	fmt.Println("submit")
	source := &models.Source{}
	path := source.GenPath()
	source.CreatedAt = time.Now()
	source.Status = models.UnHandled
	source.ProblemId = problemId
	//get user id
	username := c.Session["username"]
	u := models.GetCurrentUser(username)
	if u != nil {
		source.UserId = u.Id
	} else {
		c.Flash.Error("wrong user")
		return c.Redirect(routes.Code.Answer(problemId))
	}
	switch lang {
	case "c":
		_, err := util.WriteFile(path+"/tmp.c", []byte(code))
		if err != nil {
			fmt.Println(err)
		}
		source.Lang = models.C
	case "cpp":
		util.WriteFile(path+"/tmp.cpp", []byte(code))
		source.Lang = models.CPP
	case "go":
		util.WriteFile(path+"/tmp.go", []byte(code))
		source.Lang = models.Go
	default:
		c.Flash.Error("wrong lang %s\n", lang)
		return c.Redirect(routes.Code.Answer(problemId))
	}
	_, err := engine.Insert(source)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(routes.Code.Answer(problemId))
	}
	return c.Redirect(routes.Code.Status())
}

func (c *Code) Status() revel.Result {
	var sources []models.Source
	engine.Desc("created_at").Find(&sources)
	return c.Render(sources)
}
