package controllers

import (
	"errors"
	"fmt"
	"html/template"
	"math"
	"reflect"
	"strings"

	"github.com/revel/revel"
)

// ouput errors
type Crash struct {
	*revel.Controller
}

func (c *Crash) Notice() revel.Result {
	return c.Render()
}

/*

	This is module is not a controller but a util to show pagination
*/

type Pagination struct {
	sum     int
	current int    //current page
	url     string //template url
	pages   int    //ceiling of total pages
	hasPrev bool   //there is previous page before the current page
	hasNext bool   //there is next page after the current page
}

func (p *Pagination) isValidPage(v *revel.Validation, bean interface{}, index ...int64) {
	n, err := engine.Count(bean)
	if err != nil {
		e := &revel.ValidationError{
			Message: "bean error",
			Key:     reflect.TypeOf(bean).Name(),
		}
		v.Errors = append(v.Errors, e)
	}
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
			Message: fmt.Sprintf("%d is out of range %d", c, n/perPage),
			Key:     reflect.TypeOf(bean).Name(),
		}
		v.Errors = append(v.Errors, e)
	}
	p.current = int(c)
	if n < 0 {
		n = 0
	}
	p.sum = int(n)
}
func (p *Pagination) Page(perPage int64, url string) error {
	if p.sum < 0 {
		return errors.New("sum could not be negative number")
	}
	p.url = url
	p.hasPrev = true
	p.hasNext = true
	// ceil total number( of pages
	p.pages = int(math.Ceil(float64(p.sum) / float64(perPage)))
	if p.pages == 0 {
		p.pages = 1
	}
	if p.current < 1 || (p.current > p.pages && p.pages != 0) {
		return errors.New(fmt.Sprintf(" %d is out of range %d", p.current, p.pages))
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

func (p *Pagination) Html() template.HTML {
	html := `<div class="ui pagination menu">`
	linkFlag := "?"
	if strings.Index(p.url, "?") > -1 {
		linkFlag = "&"
	}
	if p.hasPrev {
		html += fmt.Sprintf(`<a class="icon item" href="%s%sindex=%d"><i class="icon left arrow"></i>PREV</a>`, p.url, linkFlag, p.current-1)
	}
	if p.pages == 0 {
		html += fmt.Sprintf(`<div class="disabled item">%d/%d</div>`, 0, p.pages)
	} else {
		html += fmt.Sprintf(`<div class="disabled item">%d/%d</div>`, p.current, p.pages)
	}
	if p.hasNext {
		html += fmt.Sprintf(`<a class="icon item" href="%s%sindex=%d">NEXT<i class="icon right arrow"></i></a>`, p.url, linkFlag, p.current+1)
	}
	html += `</div>`
	fmt.Println(html)
	return template.HTML(html)
}

/* replace \n with <p>*/
func Text(input string) template.HTML {
	return template.HTML("<p>" + strings.Replace(input, "\n", "<br>", -1) + "</p>")
}
