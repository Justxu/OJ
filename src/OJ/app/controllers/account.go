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
		c.Flash.Out["username"] = user.Name
		c.Flash.Error("Account or password error")
		return c.Redirect(routes.Account.Login())
	} else {
		c.Session["username"] = user.Name
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
	c.Session["username"] = user.Name
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
	username := c.Session["username"]
	user := models.GetCurrentUser(username)
	subject := "activate password"
	content := `<h2><a href="http://localhost:9000/Account/Activate/` + user.ActiveCode + `>activate account</a></h2>`
	//stmp is defined in ./init.go
	err := SendMail(subject, content, smtpConfig.Username, []string{user.Email}, smtpConfig, true)
	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Out["info"] = `<a href="/Account/ResentActiveCode">重发邮件</a>`
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

func (c Account) Forgot() revel.Result {
	return c.Render()
}

func (c Account) SendResetEmail(email string) revel.Result {
	var user models.User
	//fmt.Println(email)
	user.Email = email
	code := uuid.NewUUID()
	user.ResetCode = code.String()
	user.ResetCodeCreatedTime = time.Now()
	if user.HasEmail() {
		engine.Cols("reset_code,reset_code_created_time").Update(user)
		subject := "Reset Password"
		content := `<h2><a href="http://localhost:9000/Account/Reset/` + user.ResetCode + `">Reset Password</a></h2>`
		//stmp is defined in ./init.go
		SendMail(subject, content, smtpConfig.Username, []string{email}, smtpConfig, true)
		c.Flash.Success("Email has been sent, pleas check it.")
		return c.Redirect(routes.Account.Notice())
	} else {
		c.Flash.Error("Wrong Email")
		return c.Redirect(routes.Account.Notice())
	}
}
func (c Account) Reset(resetcode string) revel.Result {
	//fmt.Println(resetcode)
	uucode := uuid.Parse(resetcode)
	var user models.User
	has, err := engine.Where("reset_code = ?", resetcode).Get(&user)
	if err != nil {
		fmt.Println(err)
	}
	if !has {
		c.Flash.Error("wrong code")
		return c.Redirect(routes.Account.Forgot())
	}
	ut, _ := uucode.Time()
	s, n := ut.UnixTime()
	t := time.Unix(s, n)
	if user.ResetCodeCreatedTime.Sub(t) > time.Minute {
		c.Flash.Error("Reset Reply time out")
		return c.Redirect(routes.Account.Forgot())
	} else {
		c.Session["username"] = user.Name
		c.Flash.Data["resetcode"] = resetcode
		return c.Render()
	}
}

func (c Account) PostReset(user models.User) revel.Result {
	username := c.Session["username"]
	//fmt.Println("user", user)
	if user.Password == user.ConfirmPassword {
		//fmt.Println("user password ", user.Password)
		user.HashedPassword, user.Salt = models.GenHashPasswordAndSalt(user.Password)
		//fmt.Println("user Hashedpassword ", user.HashedPassword)
		//fmt.Println("user salt ", user.Salt)
		//pw := models.HashPassword("123", user.Salt)
		//fmt.Println("pw ", pw)
		user.ResetCode = ""
		//fmt.Println(username)
		_, err := engine.Where("name = ?", username).Cols("hashed_password", "salt", "reset_code").Update(&user)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println("reset ok")
		c.Session["username"] = username
		return c.Redirect(routes.Problems.Index())
	} else {
		resetcode := c.Flash.Data["resetcode"]
		//fmt.Println("post restcode", resetcode)
		c.Flash.Error("两次密码输入不一致")
		return c.Redirect(routes.Account.Reset(resetcode))
	}
}
