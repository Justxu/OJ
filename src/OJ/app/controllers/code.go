package controllers

import (
	"fmt"

	"github.com/ggaaooppeenngg/util"
	"github.com/revel/revel"

	"OJ/app/models"
	"OJ/app/routes"
)

type Code struct {
	*revel.Controller
}

func (c *Code) Index() revel.Result {
	return c.Render()
}

func (c *Code) Submit(code string) revel.Result {
	fmt.Printf("%s", code)
	fmt.Println("submit")
	fmt.Println(util.Pwd())
	source := &models.Source{}
	path := source.GenPath()
	util.WriteFile(path, []byte(code))
	out, _ := util.Run("go", "run", path)
	err := &revel.ValidationError{
		Message: string(out),
		Key:     "outErr",
	}
	c.Validation.Errors = append(c.Validation.Errors, err)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
	}
	return c.Redirect(routes.Code.Index())
}
