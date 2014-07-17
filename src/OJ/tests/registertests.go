package tests

import (
	"net/url"

	"OJ/app/models"

	"github.com/revel/revel"
)

type RegisterTest struct {
	revel.TestSuite
}

func (t *RegisterTest) Before() {
	println("Set up")
}

func (t RegisterTest) TestRegiterPageWorks() {
	t.Get("/Account/Login/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t RegisterTest) TestRegiterPostWorks() {
	var user models.User
	defer engine.Delete(&user)
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/Account/PostRegist", form)
	user.Email = "test@test.com"
	has, _ := engine.Get(&user)
	t.Assert(has)
}
func (t RegisterTest) TestActiveCode() {
	var user *models.User
	defer engine.Delete(&user)
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/Account/PostRegist", form)
	user = new(models.User)
	has, err := engine.Where("email =?", "test@test.com").Get(user)
	t.Assert(err == nil)
	t.Assert(has)
	t.Assert(user.ActiveCode != "")
	t.Get("/Account/Activate/" + user.ActiveCode)
	t.AssertOk()
	user = new(models.User)
	has, err = engine.Where("email =?", "test@test.com").Cols("active").Get(user)
	t.Assert(has)
	t.AssertEqual(err, nil)
	t.AssertEqual(user.Active, true)
}

func (t RegisterTest) TestResetCode() {
	var user *models.User
	defer engine.Delete(&user)
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/Account/PostRegist", form)
	user = new(models.User)
	has, err := engine.Where("email =?", "test@test.com").Get(user)
	t.Assert(err == nil)
	t.Assert(has)
	t.Get("/Account/Logout")
	t.PostForm("/Account/Forgot", form)
	engine.Id(user.Id).Get(user)
	t.Get("/Account/Reset/" + user.ResetCode)
	form = url.Values{
		"user.Password":        []string{"123"},
		"user.ConfirmPassword": []string{"123"},
	}
	t.PostForm("/Account/PostReset", form)
	engine.Id(user.Id).Get(user)
	pw := models.HashPassword("123", user.Salt)
	t.AssertEqual(pw, user.HashedPassword)
}

func (t *RegisterTest) After() {
	println("Tear down")
}
