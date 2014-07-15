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
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/Account/PostRegist", form)
	var user models.User
	user.Email = "test@test.com"
	has, _ := engine.Get(&user)
	t.Assert(has)
	engine.Delete(&user)
}

func (t *RegisterTest) After() {
	println("Tear down")
}
