package controllers

import (
	"fmt"

	"OJ/app/models"
	"OJ/app/routes"

	"code.google.com/p/go-uuid/uuid"
	"github.com/revel/revel"
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
	return c.Render()
}

func (c Account) SendResetEmail(email string) revel.Result {
	var user models.User
	user.Email = email
	code := uuid.NewUUID()
	user.Code = code.String()
	user.CodeCreated = code.Time()
	if user.HasEmail() {
		engine.Clos("code").Update(user)
		subject := "Reset Password"
		content := `<h2><a href="http://localhost:9000/Account/Reset/` + user.Code + `>Reset Password</a></h2>`
		SendMail(subject, body, "GOOJ", []string{email}, smtp, true)
	} else {
		c.Flash.Error("Wrong Email")
		return c.Redirect(routes.Account.Forget())
	}
}

func (c Account) Reset(code string) revel.Result {
	code := uuid.Parse(code)
	var user models.User
	engine.Where("code = ?", code).Get(&user)
	t, _ := code.Time()
	if user.CodeCreated.Sub(t) > time.Minute {
		c.Flash.Error("Reset Reply time out")
		c.Redirect(routes.Account.Forget())
	} else {
		c.Session["username"] = user.Name
		return c.Render()
	}
}

func (c Account) PostReset(user models.User) {
	username := c.Session["username"]
	if user.Password == user.ConfirmPassword {
		u.HashedPassword = models.GenHashPasswordAndSalt(u.Password)
		u.Code = ""
		engine.Id(id).Update(&user)
		c.Session["username"] = username
		c.Redirect(routes.Problems.Index())
	} else {
		c.Flash.Error("两次密码输入不一致")
		c.Redirect(routes.Account.Reset())
	}
}
