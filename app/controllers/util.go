package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"net/smtp"
	"reflect"
	"strings"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/revel/revel"
)

var (
	//用户访问权限,也就时登陆以后才可以使用的权限
	userPermission = []string{"problems.new"}
	//管理员权限
	adminPermission = []string{"problems.edit"}
	//重复登陆检查
	logoutCheck = []string{"account.login"}
)

/*

Pagination Helper

*/

type Pagination struct {
	sum     int
	current int    //current page
	url     string //template url
	pages   int    //ceiling of total pages
	hasPrev bool   //there is previous page before the current page
	hasNext bool   //there is next page after the current page
}

func (p *Pagination) isValidPage(v *revel.Validation,
	bean interface{}, index ...int64) {
	n, err := engine.Count(bean)
	if err != nil {
		e := &revel.ValidationError{
			Message: "bean error",
			Key:     reflect.TypeOf(bean).Name(),
		}
		v.Errors = append(v.Errors, e)
	}
	//modify n to 0
	if n < 0 {
		n = 0
	}
	//current page number
	var c int64
	if len(index) == 0 {
		c = 1
	} else {
		c = index[0]
		if c == 0 {
			c = 1
		}
	}
	if c*perPage > n+perPage || c < 1 {
		e := &revel.ValidationError{
			Message: fmt.Sprintf("%d is out of range %d",
				c, n/perPage),
			Key: reflect.TypeOf(bean).Name(),
		}
		v.Errors = append(v.Errors, e)
	}
}

//Page must come with is ValidPage
//Must be validated before being paged
func (p *Pagination) Page(bean interface{}, perPage int64,
	url string, index ...int64) error {
	n, _ := engine.Count(bean)
	p.sum = int(n)
	var c int64
	if len(index) == 0 {
		c = 1
	} else {
		c = index[0]
		if c == 0 {
			c = 1
		}
	}
	p.current = int(c)
	p.url = url
	p.hasPrev = true
	p.hasNext = true
	// ceil total number( of pages
	p.pages = int(math.Ceil(float64(p.sum) / float64(perPage)))
	//if sum == 0 p.pages could be 0,but it should be 1
	if p.pages == 0 {
		p.pages = 1
	}
	if p.current < 1 || (p.current > p.pages) {
		return errors.New(fmt.Sprintf(" %d is out of range %d",
			p.current, p.pages))
	} else {
		if p.current == 1 {
			p.hasPrev = false
		}
		if p.current == p.pages {
			p.hasNext = false
		}
	}
	return nil
}

//render pagination html

func (p *Pagination) Html() template.HTML {
	html := `<div class="ui pagination menu">`
	linkFlag := "/p/"
	if p.hasPrev {
		html += fmt.Sprintf(`<a class="icon item" href="%s%s%d">`+
			`<i class="icon left arrow"></i>PREV</a>`,
			p.url, linkFlag, p.current-1)
	}
	html += fmt.Sprintf(`<div class="disabled item">%d/%d</div>`,
		p.current, p.pages)
	if p.hasNext {
		html += fmt.Sprintf(`<a class="icon item" href="%s%s%d">`+
			`NEXT<i class="icon right arrow"></i></a>`,
			p.url, linkFlag, p.current+1)
	}
	html += `</div>`
	return template.HTML(html)
}

/* replace \n with <p>*/
func Text(input string) template.HTML {
	//markdown class is used by marked.js to render markdown text
	return template.HTML("<p class=\"markdown\">" + input + "</p>")
}

//checkout if user is admin
func IsAdmin(a interface{}) bool {
	if a == nil {
		return false
	} else {
		return a.(string) == admin
	}
}

//search certain content
func Search(key string, c *revel.Controller) revel.Result {
	var problems []models.Problem
	err := engine.Where("title = ? ", key).Find(&problems)
	if err != nil {
		c.Flash.Error("error %s", err.Error())
		c.Redirect(routes.Notice.Crash())
	}
	return c.Render(problems)
}

func inStringSlice(s string, slc []string) bool {
	for _, v := range slc {
		if s == v {
			return true
		}
	}
	return false
}

//check user login
func connected(c *revel.Controller) bool {
	username, has := c.Session[USERNAME]
	if !has {
		return false
	}
	u := new(models.User)
	if has, err := engine.Where("name = ?", username).
		Get(u); has && err == nil {
		return true
	} else {
		return false
	}
}

func adminAuthentication(c *revel.Controller) bool {
	username, has := c.Session[USERNAME]
	if !has {
		return false
	}
	u := new(models.User)
	if has, err := engine.Where("name = ?", username).
		Get(u); has && err == nil {
		if username == admin {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//authentication check
func authenticate(c *revel.Controller) revel.Result {
	if inStringSlice(strings.ToLower(c.Action),
		adminPermission) {
		if !adminAuthentication(c) {
			c.Flash.Error("you are not admin")
			return c.Redirect("/")
		}
	}
	if inStringSlice(strings.ToLower(c.Action),
		userPermission) {
		if ok := connected(c); !ok {
			c.Flash.Error("please login first")
			return c.Redirect(routes.Account.Login())
		} else {
			return nil
		}
	}
	if inStringSlice(strings.ToLower(c.Action),
		logoutCheck) {
		if ok := connected(c); ok {
			c.Flash.Error("can not repeat login")
			return c.Redirect("/")
		} else {
			return nil
		}
	}
	return nil
}

// SMTP util
type SmtpConfig struct {
	Username string
	Password string
	Host     string
	Addr     string
}

// send mail
func SendMail(subject string, message string, from string, to []string, smtpConfig SmtpConfig, isHtml bool) error {
	auth := smtp.PlainAuth(
		"",
		smtpConfig.Username,
		smtpConfig.Password,
		smtpConfig.Host,
	)
	contentType := "text/plain"
	if isHtml {
		contentType = "text/html"
	}
	msg := fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", strings.Join(to, ";"), from, subject, contentType, message)
	return smtp.SendMail(smtpConfig.Addr, auth, from, to, []byte(msg))
}
