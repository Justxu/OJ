package controllers

import (
	"fmt"

	"OJ/app/models"
	"OJ/app/routes"

	"github.com/ggaaooppeenngg/util"
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

func (p *Problems) PostNew(problem models.Problem, inputTest, outputTest []byte) revel.Result {
	p.Validation.Required(problem.Title)
	p.Validation.Required(problem.Description)
	p.Validation.Required(outputTest)
	p.Validation.Required(inputTest)
	fmt.Printf("out is %s\n", inputTest)
	path := problem.TestPath()
	problem.InputTest = path + "/inputTest"
	problem.OutputTest = path + "/outputTest"
	if p.Validation.HasErrors() {
		p.Validation.Keep()
		p.FlashParams()
		return p.Redirect(routes.Problems.Index())
	}
	_, err := util.WriteFile(problem.InputTest, inputTest)
	if err != nil {
		fmt.Println(err)
	}
	_, err = util.WriteFile(problem.OutputTest, outputTest)
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
func (p *Problems) Edit(problem models.Problem) revel.Result {
	engine.Update(problem)
	return p.Redirect(routes.Problems.Index())
}
func (p *Problems) Standings() revel.Result {
	return p.Redirect(routes.Problems.Index())
}
