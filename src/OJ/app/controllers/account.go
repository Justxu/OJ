package controllers

import (
	"fmt"
	"time"

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
	code := uuid.NewUUID()
	user.ActiveCode = code.String()
	user.ActiveCodeCreatedTime = time.Now()
	if !user.Save() {
		c.Flash.Error("Registered user failed")
		return c.Redirect(routes.Account.Regist())
	}
	c.Session["user"] = user.Name
	subject := "activate password"
	content := `<h2><a href="http://localhost:9000/Account/Activate/` + user.ActiveCode + `">activate account</a></h2>`
	//stmp is defined in ./init.go
	err := SendMail(subject, content, smtpConfig.Username, []string{user.Email}, smtpConfig, true)
	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Success("please check email to make your account active")
	return c.Redirect(routes.Account.Notice())
}
func (c Account) ResentActiveCode() revel.Result {
	username := c.Session["user"]
	user := models.GetCurrentUser(username)
	subject := "activate password"
	content := `<h2><a href="http://localhost:9000/Account/Activate/` + user.ActiveCode + `>activate account</a></h2>`
	//stmp is defined in ./init.go
	err := SendMail(subject, content, smtpConfig.Username, []string{user.Email}, smtpConfig, true)
	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Success("please check email to make your account active")
	return c.Redirect(routes.Account.Notice())
}

func (c Account) Activate(activecode string) revel.Result {
	var user = &models.User{}
	_, err := engine.Where("active_code = ?", activecode).Get(user)
	if err != nil {
		fmt.Println(err)
	}
	user.ActiveCode = ""
	user.Active = true
	_, err = engine.Cols("active", "active_code").Update(user)
	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Success("激活成功")
	return c.Redirect(routes.Account.Notice())
}
func (c Account) Notice() revel.Result {
	return c.Render()
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
	user.ResetCode = code.String()
	user.ResetCodeCreatedTime = time.Now()
	if user.HasEmail() {
		engine.Cols("reset_code,reset_code_created_time").Update(user)
		subject := "Reset Password"
		content := `<h2><a href="http://localhost:9000/Account/Reset/` + user.ResetCode + `>Reset Password</a></h2>`
		//stmp is defined in ./init.go
		SendMail(subject, content, "GOOJ", []string{email}, smtpConfig, true)
		return c.Redirect(routes.Account.Forget())
	} else {
		c.Flash.Error("Wrong Email")
		return c.Redirect(routes.Account.Forget())
	}
}
func (c Account) Forgot() revel.Result {
	return c.Render()
}
func (c Account) Reset(code string) revel.Result {
	uucode := uuid.Parse(code)
	var user models.User
	engine.Where("code = ?", code).Get(&user)
	ut, _ := uucode.Time()
	s, n := ut.UnixTime()
	t := time.Unix(s, n)
	if user.ResetCodeCreatedTime.Sub(t) > time.Minute {
		c.Flash.Error("Reset Reply time out")
		return c.Redirect(routes.Account.Forget())
	} else {
		c.Session["username"] = user.Name
		return c.Render()
	}
}

func (c Account) PostReset(user models.User) revel.Result {
	username := c.Session["username"]
	if user.Password == user.ConfirmPassword {
		user.HashedPassword, user.Salt = models.GenHashPasswordAndSalt(user.Password)
		user.ResetCode = ""
		engine.Where("username = ?", username).Update(&user)
		c.Session["username"] = username
		return c.Redirect(routes.Problems.Index())
	} else {
		c.Flash.Error("两次密码输入不一致")
		return c.Redirect(routes.Account.Reset(user.ResetCode))
	}
}
