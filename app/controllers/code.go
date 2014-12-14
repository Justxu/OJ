package controllers

import (
	"log"
	"time"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"
)

const (
	ERROR  = "error"
	PANIC  = "panic"
	CODE   = "code"
	STATUS = "status"
	REPORT = "report"
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
	source := &models.Source{}
	path := source.GenPath()
	source.CreatedAt = time.Now()
	source.Status = models.UnHandled
	source.ProblemId = problemId
	//get user id
	username := c.Session[USERNAME]
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
			log.Println(err)
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
	return c.Redirect("/code/status")
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
		return c.Redirect(routes.Notice.Crash())
	}
	err := engine.Desc("created_at").
		Asc("id").
		Limit(perPage, perPage*(pagination.current-1)).
		Find(&sources)
	if err != nil {
		log.Println(err)
	}
	err = pagination.Page(models.Source{},
		perPage,
		"/code/status",
		index)
	if err != nil {
		c.Flash.Error("pagination error")
		c.Redirect(routes.Notice.Crash())
	}
	moreScripts = append(moreStyles,
		"js/prettify.js",
		"js/code_status.js")
	moreStyles = append(moreStyles, "css/prettify.css")
	return c.Render(moreStyles, moreScripts, sources, pagination)
}

//json render

//get source code
func (c *Code) View(id int64) revel.Result {
	s := models.Source{}
	has, err := engine.Id(id).Get(&s)
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})
	if !has {
		data[STATUS] = false
		data[ERROR] = "not exits"
		return c.RenderJson(data)
	}
	code, err := s.View()
	if err != nil {
		data[STATUS] = false
		data[ERROR] = err.Error()
	}
	data[STATUS] = true
	data[CODE] = code
	return c.RenderJson(data)
}

//get panic error of the code
func (c *Code) GetPanic(id int64) revel.Result {
	s := new(models.Source)
	has, err := engine.Id(id).Get(s)
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})
	if !has {
		data[ERROR] = "not exist!"
		data[STATUS] = false
	} else {
		data[PANIC] = s.PanicError
		data[STATUS] = true
	}
	return c.RenderJson(data)
}

// check the output
func (c *Code) Check(id int64) revel.Result {
	s := models.Source{}
	has, err := engine.Id(id).Get(&s)
	if err != nil {
		log.Println(err)
	}
	data := make(map[string]interface{})
	if !has {
		data[STATUS] = false
		data[ERROR] = "not exist!"
		return c.RenderJson(data)
	}
	r, e := s.Check()
	if e != nil || s.Status != models.WrongAnswer {
		data[STATUS] = false
		if e != nil {
			data[ERROR] = e.Error()
		} else {
			data[ERROR] = "not wrong answer"
		}
		return c.RenderJson(data)
	} else {
		data[STATUS] = true
		data[REPORT] = r
		return c.RenderJson(data)
	}
}
