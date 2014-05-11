package controllers

import (
	"fmt"

	"github.com/ggaaooppeenngg/OJ/app/routes"
	"github.com/ggaaooppeenngg/util"

	"github.com/revel/revel"
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
	util.WriteFile("github.com/ggaaooppeenngg/OJ/tmp/main.go", []byte(code))
	out, _ := util.Run("go", "run", "github.com/ggaaooppeenngg/OJ/tmp/main.go")
	err := &revel.ValidationError{
		Message: string(out),
		Key:     "outErr",
	}
	fmt.Println(string(out))
	c.Validation.Errors = append(c.Validation.Errors, err)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
	}
	return c.Redirect(routes.Code.Index())
}
