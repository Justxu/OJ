package tests

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/ggaaooppeenngg/OJ/app/models"

	"github.com/revel/revel"
)

type RegisterTest struct {
	revel.TestSuite
}

func (t *RegisterTest) Before() {
	cookieJar, _ := cookiejar.New(nil)
	t.Client = &http.Client{
		Jar: cookieJar,
	}
}

func (t *RegisterTest) TestRegisterPageWorks() {
	t.Get("/account/login/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *RegisterTest) TestRegisterPostWorks() {
	var user models.User
	defer engine.Delete(&user)
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/account/postregister", form)
	user.Email = "test@test.com"
	has, _ := engine.Get(&user)
	t.Assert(has)
}

func (t *RegisterTest) TestActiveCode() {
	var user models.User
	defer engine.Delete(&user)
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/account/register", form)
	has, err := engine.Where("email =?", "test@test.com").Get(&user)
	t.Assert(err == nil)
	t.Assert(has)
	t.Assert(user.ActiveCode != "")
	t.Get("/account/activate/" + user.ActiveCode)
	fmt.Println(user.Id)
	user = models.User{}
	has, err = engine.Where("email =?", "test@test.com").Get(&user)
	fmt.Println(user)
	t.AssertEqual(nil, err)
	t.Assert(has)
	t.AssertEqual(true, user.Active)
}

func (t *RegisterTest) TestResetCode() {
	var user models.User
	defer engine.Delete(&user)
	form := url.Values{
		"user.Name":            []string{"testName"},
		"user.Email":           []string{"test@test.com"},
		"user.Password":        []string{"testtest"},
		"user.ConfirmPassword": []string{"testtest"},
	}
	t.PostForm("/account/register", form)
	has, err := engine.Where("email =?", "test@test.com").Get(&user)
	t.Assert(err == nil)
	t.Assert(has)
	username := t.Session["username"]
	t.AssertEqual("testName", username)
	t.Get("/account/logout")
	username, has = t.Session["username"]
	t.Assert(!has)
	form = url.Values{
		"email": []string{"test@test.com"},
	}
	t.PostForm("/account/send-reset-email", form)
	user = models.User{}
	has, _ = engine.Where("email =?", "test@test.com").Get(&user)
	t.Assert(has)
	t.Assert(user.ResetCode != "")
	t.Get("/account/reset/" + user.ResetCode)
	username = t.Session["username"]
	t.AssertEqual("testName", username)
	form = url.Values{
		"user.Password":        []string{"123"},
		"user.ConfirmPassword": []string{"123"},
	}
	t.PostForm("/account/reset", form)
	user = models.User{}
	has, _ = engine.Where("email =?", "test@test.com").Get(&user)
	t.Assert(has)
	//println(user.Salt)
	pw := models.HashPassword("123", user.Salt)
	t.AssertEqual(user.HashedPassword, pw)
}

func (t *RegisterTest) After() {
	println("Tear down")
}
