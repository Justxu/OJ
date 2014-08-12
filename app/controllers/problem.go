package controllers

import (
	"fmt"
	"path"
	"strconv"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"
)

type Problems struct {
	*revel.Controller
}

func (p *Problems) Index() revel.Result {
	var problems []models.Problem
	err := engine.Limit(10).Find(&problems)
	if err != nil {
		fmt.Println(err)
	}
	return p.Render(problems)
}

//URL: prolem/p/id,get problem information
func (p *Problems) P(index int) revel.Result {
	p.Validation.Min(index, 0).Message("worong problem index")
	if p.Validation.HasErrors() {
		p.FlashParams()
		p.Validation.Keep()
		return p.Redirect(routes.Problems.Index())
	}
	var prob models.Problem
	has, err := engine.Id(index).Get(&prob)
	if err != nil || !has {
		fmt.Println(err)
		p.Flash.Error("problem id error %d", index)
		p.Redirect(routes.Problems.Index())
	}
	return p.Render(prob)
}

func (p *Problems) PostNew(problem models.Problem, inputTest, outputTest []byte) revel.Result {
	p.Validation.Required(problem.Title)
	p.Validation.Min(int(problem.MemoryLimit), 1).Message("TimeLimit Required")
	p.Validation.Min(int(problem.TimeLimit), 1).Message("MemoryLimit Required")
	p.Validation.Required(problem.Description)
	p.Validation.Required(outputTest).Message("output file needed")
	p.Validation.Required(inputTest).Message("input file needed")
	fmt.Printf("out is %s\n", inputTest)
	path := problem.TestPath()
	problem.InputTestPath = path + "/inputTest"
	problem.OutputTestPath = path + "/outputTest"
	if p.Validation.HasErrors() {
		p.Validation.Keep()
		p.FlashParams()
		return p.Redirect(routes.Problems.Index())
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
	return p.Redirect(routes.Problems.Index())
}

func (p *Problems) New() revel.Result {
	return p.Render()
}

func (p *Problems) Delete(id int64) revel.Result {
	problem := &models.Problem{Id: id}
	engine.Delete(problem)
	return p.Redirect(routes.Problems.Index())
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
	}
	if outputTest != nil {
		problem.OutputTestPath = path.Dir(problem.OutputTestPath) + "/outputTest"
	}
	fmt.Println("update id is")
	id, err := strconv.ParseInt(p.Session["id"], 10, 64)
	if err != nil {
		p.Flash.Error("id error")
		fmt.Println(err)
		return p.Redirect(routes.Problems.Index())
	}
	_, err = engine.Id(id).Update(problem)
	if err != nil {
		fmt.Println(err)
	}
	return p.Redirect(routes.Problems.P(int(id)))
}

func (p *Problems) Standings() revel.Result {
	return p.Redirect(routes.Problems.Index())
}
