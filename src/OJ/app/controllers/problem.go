package controllers

import (
	"bufio"
	"fmt"
	"os"

	"OJ/app/models"
	"OJ/app/routes"

	"github.com/revel/revel"
)

type Problems struct {
	*revel.Controller
}

func (p Problems) Index() revel.Result {
	var problems []models.Problem
	err := engine.Limit(10).Find(&problems)
	if err != nil {
		fmt.Println(err)
	}
	return p.Render(problems)
}

func (p *Problems) PostNew(problem models.Problem, inputTest, outputTest io.Reader) revel.Result {
	p.Validation.Required(problem.Title)
	p.Validation.Required(problem.Description)
	p.Validation.Required(outputTest)
	p.Validation.Required(inputTest)
	pathIn := problem.TestPath("outputTest")
	pathOut := problem.TestPath("inputTest")
	defer inputTest.Close()
	defer outputTest.Close()
	if p.Validation.HasErrors() {
		return p.Redirect(routes.Problems.Index())
	}

	_, err := util.CopyFile(pathIn, inputTest)
	_, err = util.CopyFile(pathOut, outputTest)
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
func (p *Problems) Edit(problem models.Problem) revel.Result {
	engine.Update(problem)
	return p.Redirect(routes.Problems.Index())
}
