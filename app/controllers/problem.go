package controllers

import (
	"fmt"
	"log"
	"path"
	"strconv"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"
)

const (
	ID      = "id"
	perPage = 5
)

type Problem struct {
	*revel.Controller
}

//GET problem/p/:index,get problem information
func (p *Problem) Index(index int64) revel.Result {
	var problems []models.Problem
	pagination := &Pagination{}
	pagination.isValidPage(p.Validation, models.Problem{}, index)

	if p.Validation.HasErrors() {
		p.FlashParams()
		p.Validation.Keep()
		return p.Redirect(routes.Notice.Crash())
	}
	err := pagination.Page(models.Problem{},
		perPage,
		"/problem",
		index)

	if err != nil {
		p.Flash.Error(err.Error())
		p.Redirect(routes.Notice.Crash())
	}

	err = engine.Asc(ID).
		Where("is_valid = ?", true).
		Limit(perPage, perPage*(pagination.current-1)).
		Find(&problems)
	if err != nil {
		log.Println(err)
	}

	return p.Render(problems, pagination)
}

//GET problem/:i,get problem information
func (p *Problem) P(id int) revel.Result {
	p.Validation.Min(id, 0).Message("worong problem i")
	if p.Validation.HasErrors() {
		//把参数存到Flash里面方便debug
		//p.FlashParams()
		p.Validation.Keep()
		return p.Redirect("/")
	}
	var prob models.Problem
	has, err := engine.
		Id(id).
		Get(&prob)

	if err != nil || !has || !prob.IsValid {
		log.Println(err)
		p.Flash.Error("problem id error %d", id)
		p.Redirect("/")
	}
	//markdown plugin script
	var moreScripts []string
	moreScripts = append(moreScripts, "js/marked.js")
	return p.Render(prob, moreScripts)
}

//POST /problem/new
func (p *Problem) PostNew(problem models.Problem,
	inputTest, outputTest []byte) revel.Result {

	problem.Validate(p.Validation, inputTest, outputTest)
	if p.Validation.HasErrors() {
		p.Validation.Keep()
		//p.FlashParams()
		return p.Redirect("/")
	}
	if IsAdmin(p.Session[USERNAME]) {
		problem.IsValid = true
	} else {
		//if user is not admin,checked
		//the problem effectiveness manually by administrators
		has, id := models.GetUserId(p.Session[USERNAME])
		if has {
			problem.PosterId = id
		}
		problem.IsValid = false
	}
	_, err := util.WriteFile(problem.InputTestPath, inputTest)
	if err != nil {
		p.Flash.Error(err.Error())
		log.Println(err)
		return p.Redirect(routes.Notice.Crash())
	}
	_, err = util.WriteFile(problem.OutputTestPath, outputTest)
	if err != nil {
		p.Flash.Error(err.Error())
		log.Println(err)
		return p.Redirect(routes.Notice.Crash())
	}
	_, err = engine.Insert(&problem)
	if err != nil {
		p.Flash.Error("insert error")
		log.Println(err)
		return p.Redirect(routes.Notice.Crash())
	}
	p.Flash.Success("post success!")
	return p.Redirect("/")
}

//list unchecked users' problem posts
//GET /problem/posts/p/:index
func (p *Problem) Posts(index int64) revel.Result {
	var problems []models.Problem
	pagination := &Pagination{}
	pagination.isValidPage(p.Validation, models.Problem{}, index)
	if p.Validation.HasErrors() {
		//p.FlashParams()
		p.Validation.Keep()
		return p.Redirect(routes.Notice.Crash())
	}
	err := engine.Asc(ID).
		Where("is_valid = ?", false).
		Limit(perPage, perPage*(pagination.current-1)).
		Find(&problems)
	if err != nil {
		fmt.Println(err)
	}
	err = pagination.Page(models.Problem{}, perPage,
		p.Request.Request.URL.Path, index)
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	return p.Render(problems)
}

//GET /problem/admin/:id , make problem valid
func (p *Problem) Admit(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	problem.IsValid = true
	_, err := engine.Cols("is_valid").Update(problem)
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	return p.Redirect("/")
}

func (p *Problem) New() revel.Result {
	return p.Render()
}

// GET /problem/delete/:id
func (p *Problem) Delete(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	has, err := engine.Get(problem)
	if !has {
		p.Flash.Error("problem invalid")
		return p.Redirect(routes.Notice.Crash())
	}
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	err = problem.Delete()
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	return p.Redirect("/")
}

// GET /problem/edit/:id
func (p *Problem) Edit(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	has, err := engine.Id(problem.Id).Get(problem)
	if err != nil || !has {
		p.Flash.Error("no such problem")
		return p.Redirect(routes.Notice.Crash())
	}
	p.Session[ID] = strconv.Itoa(int(id))
	return p.Render(problem)
}

// POST /problem/edit
func (p *Problem) PostEdit(problem models.Problem,
	inputTest, outputTest []byte) revel.Result {
	defer func() {
		delete(p.Session, ID)
	}()
	if inputTest != nil {
		problem.InputTestPath = path.Dir(problem.InputTestPath) +
			"/inputTest"
		_, err := util.WriteFile(problem.InputTestPath, inputTest)
		if err != nil {
			log.Println(err)
		}

	}
	if outputTest != nil {
		problem.OutputTestPath = path.Dir(problem.OutputTestPath) +
			"/outputTest"
		_, err := util.WriteFile(problem.OutputTestPath, outputTest)
		if err != nil {
			log.Println(err)
		}
	}
	id, err := strconv.ParseInt(p.Session[ID], 10, 64)
	if err != nil {
		p.Flash.Error("id error")
		log.Println(err)
		return p.Redirect("/")
	}
	_, err = engine.Id(id).Update(problem)
	if err != nil {
		log.Println(err)
	}
	return p.Redirect("/")
}
