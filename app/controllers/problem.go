package controllers

import (
	"fmt"
	"path"
	"regexp"
	"strconv"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"
)

const (
	perPage = 5
)

type Problems struct {
	*revel.Controller
}

var (
	urlMatch *regexp.Regexp
)

func validPicture(url string) bool {
	return true
}

//URL: prolem/Index/index,get problem information
func (p *Problems) Index(index int64) revel.Result {
	var problems []models.Problem
	pagination := &Pagination{}
	pagination.isValidPage(p.Validation, models.Problem{}, index)
	if p.Validation.HasErrors() {
		p.FlashParams()
		p.Validation.Keep()
		return p.Redirect(routes.Notice.Crash())
	}

	err := engine.Asc("id").Where("is_valid = ?", true).Limit(perPage, perPage*(pagination.current-1)).Find(&problems)
	if err != nil {
		fmt.Println(err)
	}

	err = pagination.Page(perPage, p.Request.Request.URL.Path)
	if err != nil {
		p.Flash.Error(err.Error())
		p.Redirect(routes.Notice.Crash())
	}
	return p.Render(problems, pagination)
}

//URL: prolem/p/id,get problem information
func (p *Problems) P(index int) revel.Result {
	p.Validation.Min(index, 0).Message("worong problem index")
	if p.Validation.HasErrors() {
		p.FlashParams()
		p.Validation.Keep()
		return p.Redirect(routes.Problems.Index(0))
	}
	var prob models.Problem
	has, err := engine.Id(index).Get(&prob)
	if err != nil || !has || !prob.IsValid {
		fmt.Println(err)
		p.Flash.Error("problem id error %d", index)
		p.Redirect(routes.Problems.Index(0))
	}
	//markdown script
	var moreScripts []string
	moreScripts = append(moreScripts, "js/marked.js")
	return p.Render(prob, moreScripts)
}

func (p *Problems) PostNew(problem models.Problem, inputTest, outputTest []byte) revel.Result {
	p.Validation.Required(problem.Title).Message("Title Required")
	p.Validation.Min(int(problem.MemoryLimit), 1).Message("TimeLimit Required")
	p.Validation.Min(int(problem.TimeLimit), 1).Message("MemoryLimit Required")
	p.Validation.Required(problem.Description).Message("Description Required")
	p.Validation.Required(outputTest).Message("output file needed")
	p.Validation.Required(inputTest).Message("input file needed")
	p.Validation.MaxSize(problem.InputSample, 512).Message("input sample too long")
	p.Validation.MaxSize(problem.OutputSample, 512).Message("output sample too long")

	path := problem.TestPath()
	problem.InputTestPath = path + "/inputTest"
	problem.OutputTestPath = path + "/outputTest"
	if p.Validation.HasErrors() {
		p.Validation.Keep()
		p.FlashParams()
		return p.Redirect(routes.Problems.Index(0))
	}
	if IsAdmin(p.Session["username"]) {
		problem.IsValid = true
	} else {
		//if user is not admin,checked the problem effectiveness manually by administrators
		problem.IsValid = false
	}
	_, err := util.WriteFile(problem.InputTestPath, inputTest)
	if err != nil {
		fmt.Println(err)
	}
	_, err = util.WriteFile(problem.OutputTestPath, outputTest)
	if err != nil {
		fmt.Println(err)
	}
	_, err = engine.Insert(&problem)
	if err != nil {
		fmt.Print(err)
	}
	return p.Redirect(routes.Problems.Index(0))
}

//list unchecked users' problem posts
func (p *Problems) Posts(index int64) revel.Result {
	var problems []models.Problem
	pagination := &Pagination{}
	pagination.isValidPage(p.Validation, models.Problem{}, index)
	if p.Validation.HasErrors() {
		p.FlashParams()
		p.Validation.Keep()
		return p.Redirect(routes.Notice.Crash())
	}
	err := engine.Asc("id").Where("is_valid = ?", false).Limit(perPage, perPage*(pagination.current-1)).Find(&problems)
	if err != nil {
		fmt.Println(err)
	}
	err = pagination.Page(perPage, p.Request.Request.URL.Path)
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	return p.Render(problems)
}
func (p *Problems) Admit(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	problem.IsValid = true
	_, err := engine.Cols("is_valid").Update(problem)
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	return p.Redirect(routes.Problems.Index(0))
}

func (p *Problems) New() revel.Result {
	return p.Render()
}

func (p *Problems) Delete(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	_, err := engine.Delete(problem)
	if err != nil {
		p.Flash.Error(err.Error())
		return p.Redirect(routes.Notice.Crash())
	}
	return p.Redirect(routes.Problems.Index(0))
}

func (p *Problems) Edit(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	engine.Id(problem.Id).Get(problem)
	p.Session["id"] = strconv.Itoa(int(id))
	fmt.Println("eidt id is", strconv.Itoa(int(id)))
	return p.Render(problem)
}

func (p *Problems) EditPost(problem models.Problem, inputTest, outputTest []byte) revel.Result {
	defer func() {
		delete(p.Session, "id")
	}()
	if inputTest != nil {
		problem.InputTestPath = path.Dir(problem.InputTestPath) + "/inputTest"
		_, err := util.WriteFile(problem.InputTestPath, inputTest)
		if err != nil {
			fmt.Println(err)
		}

	}
	if outputTest != nil {
		problem.OutputTestPath = path.Dir(problem.OutputTestPath) + "/outputTest"
		_, err := util.WriteFile(problem.OutputTestPath, outputTest)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("update id is")
	id, err := strconv.ParseInt(p.Session["id"], 10, 64)
	if err != nil {
		p.Flash.Error("id error")
		fmt.Println(err)
		return p.Redirect(routes.Problems.Index(0))
	}
	_, err = engine.Id(id).Update(problem)
	if err != nil {
		fmt.Println(err)
	}
	return p.Redirect(routes.Problems.P(int(id)))
}

func (p *Problems) Standings() revel.Result {
	return p.Redirect(routes.Problems.Index(0))
}
