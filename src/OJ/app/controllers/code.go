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

func (c *Code) Submit(code string, problemId int64) revel.Result {
	fmt.Println("submit")
	source := &models.Source{}
	path := source.GenPath()
	source.CreatedAt = time.Now()
	source.Status = models.UnHandled
	//
	source.ProblemId = problemId
	//我自己
	has, id := models.GetUserId(c.Session["user"])
	if has {
		source.UserId = id
	} else {
		c.Flash.Error("error")
		return c.Redirect(routes.Code.Answer(problemId))
	}
	util.WriteFile(path, []byte(code))
	engine.Insert(source)
	return c.Redirect(routes.Code.Status())
}

func (c *Code) Status() revel.Result {
	var sources []models.Source
	engine.Desc("created_at").Find(&sources)
	return c.Render(sources)
}
