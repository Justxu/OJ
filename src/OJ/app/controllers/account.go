package controllers

import (
	"fmt"

	"github.com/revel/revel"

	"OJ/app/models"
	"OJ/app/routes"
)

type Account struct {
	*revel.Controller
}

func (c Account) Login() revel.Result {
	fmt.Println("login")
	return c.Render()
}

func (c Account) PostLogin(user models.User) revel.Result {
	c.Validation.Email(user.Email)
	c.Validation.Required(user.Password)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Account.Login())
	}
	if !user.LoginOk() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Out["user"] = user.Name
		c.Flash.Error("Account or password error")
		return c.Redirect(routes.Account.Login())
	}
	return c.Redirect(routes.Code.Status())
}

func (c *Account) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.Problems.Index())
}

func (c Account) PostRegist(user models.User) revel.Result {
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		//把错误信息存到flash
		c.Validation.Keep()
		//把参数存到flash
		c.FlashParams()
		return c.Redirect(routes.Account.Regist())
	}

	if !user.Save() {
		c.Flash.Error("Registered user failed")
		return c.Redirect(routes.Account.Regist())
	}
	c.Session["user"] = user.Name
	return c.Redirect(routes.Account.Regist())
}
func (c Account) Regist() revel.Result {
	return c.Render()
}

func (c Account) Forget() revel.Result {
	return c.Redirect(routes.Account.Login())
}

func (c Account) Reset() revel.Result {
	return c.Redirect(routes.Account.Login())
}
