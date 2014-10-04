package controllers

import (
	"fmt"
	"time"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"
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
	c.Flash.Success(username + " 提交成功")
	return c.Redirect(routes.Code.Status(0))
}

func (c *Code) Status(index int64) revel.Result {
	var moreStyles []string
	var moreScripts []string
	var sources []models.Source
	pagination := &Pagination{}
	pagination.isValidPage(c.Validation, models.Source{}, index)
	if c.Validation.HasErrors() {
		c.FlashParams()
		c.Validation.Keep()
		return c.Redirect(routes.Crash.Notice())
	}
	err := engine.Desc("created_at").Limit(perPage, perPage*(pagination.current-1)).Find(&sources)
	if err != nil {
		fmt.Println(err)
	}
	err = pagination.Page(perPage, c.Request.Request.URL.Path)
	if err != nil {
		c.Flash.Error("pagination error")
		c.Redirect(routes.Crash.Notice())
	}
	moreScripts = append(moreStyles, "js/prettify.js", "js/code_status.js")
	moreStyles = append(moreStyles, "css/prettify.css")
	return c.Render(moreStyles, moreScripts, sources, pagination)
}

//get source code
func (c *Code) View(id int64) revel.Result {
	s := models.Source{}
	has, _ := engine.Id(id).Get(&s)
	data := make(map[string]interface{})
	if !has {
		data["status"] = false
		data["error"] = "not exits"
		return c.RenderJson(data)
	}
	code, err := s.View()
	if err != nil {
		data["status"] = false
		data["error"] = err.Error()
	}
	data["status"] = true
	data["code"] = code
	return c.RenderJson(data)
}

func (c *Code) Check(id int64) revel.Result {
	s := models.Source{}
	has, _ := engine.Id(id).Get(&s)
	data := make(map[string]interface{})
	if !has {
		data["status"] = false
		data["error"] = "not exist!"
		return c.RenderJson(data)
	}
	r, e := s.Check()
	if e != nil || s.Status != models.WrongAnswer {
		data["status"] = false
		if e != nil {
			data["error"] = e.Error()
		} else {
			data["error"] = "not wrong answer"
		}
		return c.RenderJson(data)
	} else {
		data["status"] = true
		data["report"] = r
		return c.RenderJson(data)
	}
}
