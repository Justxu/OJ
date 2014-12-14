package controllers

import (
	"fmt"
	"time"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"code.google.com/p/go-uuid/uuid"
	"github.com/ftrvxmtrx/gravatar"
	"github.com/revel/revel"
)

const (
	USERNAME  = "username"
	RESETCODE = "resetcode"
)

type controller interface {
}

type Account struct {
	*revel.Controller
}

//GET /account/login
//TODO: chang to general render
func (c Account) Login() revel.Result {
	return c.Render()
}

//POST /account/login
func (c *Account) PostLogin(user models.User) revel.Result {
	//check input legality
	c.Validation.Email(user.Email)
	c.Validation.Required(user.Password)
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		//c.FlashParams()
		return c.Redirect(routes.Account.Login())
	}

	if !user.LoginOk() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("account or password is incorrect")
		return c.Redirect(routes.Account.Login())
	} else {
		c.Session[USERNAME] = user.Name
	}
	return c.Redirect("/code/status")
}

//GET /account/logout
func (c *Account) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect("/")
}

//generate gravatar url
func genGravatarUrl(email string) string {
	return gravatar.GetAvatarURL(
		"https",
		gravatar.EmailHash(email),
		gravatar.DefaultMonster,
		256).String()

}

//POST account/register
func (c Account) PostRegister(user models.User) revel.Result {
	user.Validate(c.Validation)
	if c.Validation.HasErrors() {
		//把错误信息存到flash
		c.Validation.Keep()
		//把参数存到flash
		//c.FlashParams()
		return c.Redirect(routes.Account.Register())
	}
	user.GravatarUrl = genGravatarUrl(user.Email)
	code := uuid.NewUUID()
	user.ActiveCode = code.String()
	user.ActiveCodeCreatedTime = time.Now()
	if !user.Save() {
		c.Flash.Error("Registered user failed")
		return c.Redirect(routes.Account.Register())
	}
	c.Session[USERNAME] = user.Name
	subject := "activate password"
	content := fmt.Sprintf(`<h2><a href="http://%s:%s/account/activate/%s">`+
		`activate account</a></h2>`,
		appAddr, appPort, user.ActiveCode)
	err := SendMail(
		subject,
		content,
		smtpConfig.Username,
		[]string{user.Email},
		smtpConfig,
		true)
	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Success("please check email to make your account active")
	return c.Redirect(routes.Account.Notice())
}

//GET /account/resent-active-code
func (c Account) ResentActiveCode() revel.Result {
	username := c.Session[USERNAME]
	user := models.GetCurrentUser(username)
	subject := "activate password"
	content := fmt.Sprintf(`<h2><a href="http://%s:%s/account/activate/%s">`+
		`activate account</a></h2>`,
		appAddr, appPort, user.ActiveCode)
	err := SendMail(
		subject,
		content,
		smtpConfig.Username,
		[]string{user.Email},
		smtpConfig,
		true)

	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Success("please check email " +
		user.Email +
		"to make your account active")
	return c.Redirect(routes.Account.Notice())
}

//GET /account/activate/:activecode
func (c Account) Activate(activecode string) revel.Result {
	var user = &models.User{}
	has, err := engine.Where("active_code = ?", activecode).Get(user)
	if !has {
		c.Flash.Error("incorrect active code")
		return c.Redirect(routes.Account.Notice())
	}
	if err != nil {
		fmt.Println(err)
	}
	user.ActiveCode = ""
	user.Active = true
	_, err = engine.Cols("active", "active_code").Update(user)
	if err != nil {
		fmt.Println(err)
	}
	c.Flash.Success("activated!")
	return c.Redirect(routes.Problem.Index(0))
}

//GET /account/notice
func (c Account) Notice() revel.Result {
	return c.Render()
}
func (c Account) Register() revel.Result {
	return c.Render()
}

func (c Account) Forgot() revel.Result {
	return c.Render()
}

//POST /account/send-reset-email
func (c Account) SendResetEmail(email string) revel.Result {
	var user models.User
	code := uuid.NewUUID()
	user.Email = email
	user.ResetCode = code.String()
	user.ResetCodeCreatedTime = time.Now()
	if user.HasEmail() {
		_, err := engine.Where("email = ?", email).
			Cols("reset_code", "reset_code_created_time").
			Update(&user)
		if err != nil {
			fmt.Println(err)
		}
		subject := "reset password"
		content := fmt.Sprintf(
			`<h2><a href="http://%s:%s/account/reset/%s">Reset Password</a></h2>`,
			appAddr, appPort, user.ResetCode)
		SendMail(
			subject,
			content,
			smtpConfig.Username,
			[]string{email},
			smtpConfig,
			true)
		c.Flash.Success("Email has been sent, pleas check it.")
		return c.Redirect(routes.Account.Notice())
	} else {
		c.Flash.Error("Incorrect Email")
		return c.Redirect(routes.Account.Notice())
	}
}

//get time.Time form uuid code
func getTime(uucode uuid.UUID) time.Time {
	ut, _ := uucode.Time()
	s, n := ut.UnixTime()
	return time.Unix(s, n)
}

//GET /account/reset/:resetcode
func (c Account) Reset(resetcode string) revel.Result {
	uucode := uuid.Parse(resetcode)
	var user models.User
	has, err := engine.
		Where("reset_code = ?", resetcode).
		Get(&user)
	if err != nil {
		fmt.Println(err)
	}
	if !has {
		c.Flash.Error("wrong code")
		return c.Redirect(routes.Account.Forgot())
	}
	if user.ResetCodeCreatedTime.Sub(getTime(uucode)) >
		10*time.Minute {

		c.Flash.Error("Reset Reply time out")
		return c.Redirect(routes.Account.Forgot())
	} else {
		c.Session[USERNAME] = user.Name
		c.Flash.Data[RESETCODE] = resetcode
		return c.Render()
	}
}

//POST /account/reset
func (c Account) PostReset(user models.User) revel.Result {
	username := c.Session[USERNAME]
	if user.Password == user.ConfirmPassword {
		user.HashedPassword, user.Salt = models.GenHashPasswordAndSalt(user.Password)
		user.ResetCode = ""
		_, err := engine.Where("name = ?", username).
			Cols("hashed_password", "salt", "reset_code").
			Update(&user)
		if err != nil {
			fmt.Println(err)
		}
		c.Session[USERNAME] = username
		c.Flash.Success("reset success!")
		return c.Redirect(routes.Problem.Index(0))
	} else {
		resetcode := c.Flash.Data[RESETCODE]
		c.Flash.Error("两次密码输入不一致")
		return c.Redirect(routes.Account.Reset(resetcode))
	}
}

//GET /account/edit
func (c Account) Edit() revel.Result {
	username := c.Session[USERNAME]
	return c.Render(username)
}

//POST /account/edit
func (c Account) PostEdit(user models.User) revel.Result {
	c.Validation.Required(user.Name).Message("用户名不能为空")
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Account.Edit())
	}
	if user.HasName() {
		return c.Redirect(routes.User.Profile())
	}
	if user.Password != "" {
		if user.Password == user.ConfirmPassword {
			user.HashedPassword, user.Salt = models.GenHashPasswordAndSalt(user.Password)
			username := c.Session[USERNAME]
			u := models.GetCurrentUser(username)
			_, err := engine.Id(u.Id).Update(user)
			if err != nil {
				fmt.Println(err)
			}
			c.Session[USERNAME] = user.Name
			return c.Redirect(routes.User.Profile())
		} else {
			c.Flash.Error("passwords not match")
			return c.Redirect(routes.Account.Notice())
		}
	} else {
		username := c.Session[USERNAME]
		u := models.GetCurrentUser(username)
		_, err := engine.Id(u.Id).Cols("name").Update(user)
		if err != nil {
			fmt.Println(err)
			c.Flash.Error(err.Error())
		} else {
			c.Session[USERNAME] = user.Name
			c.Flash.Success("modify sucess")
		}
		return c.Redirect(routes.User.Profile())
	}

}
