package controllers

import (
	"log"

	"github.com/ggaaooppeenngg/OJ/app/models"
	"github.com/ggaaooppeenngg/OJ/app/routes"

	"github.com/revel/revel"
)

type User struct {
	*revel.Controller
}

//URL: /user/rating/p/:index
func (u *User) Rating(index int64) revel.Result {
	var users []models.User
	pagination := &Pagination{}
	pagination.isValidPage(u.Validation, models.User{}, index)
	if u.Validation.HasErrors() {
		//u.FlashParams()
		u.Validation.Keep()
		return u.Redirect(routes.Notice.Crash())
	}
	err := engine.Limit(perPage, perPage*(pagination.current-1)).
		Desc("solved").
		Asc("id").
		Find(&users)
	if err != nil {
		log.Println(err)
	}
	err = pagination.Page(models.User{}, perPage, "/user/rating", index)
	if err != nil {
		u.Flash.Error(err.Error())
		u.Redirect(routes.Notice.Crash())
	}
	return u.Render(users, pagination)
}

//GET /user/u/:id
func (u *User) ProfileVisit(id int64) revel.Result {
	var user models.User
	has, err := engine.Id(id).Get(&user)
	if !has || err != nil {
		u.Redirect(routes.Notice.Crash())
	}
	return u.Render(user)
}

//GET /user/profile
func (u *User) Profile() revel.Result {
	var user models.User
	username := u.Session[USERNAME]
	engine.Where("name = ?", username).Get(&user)
	return u.Render(user)
}

//GET /user/solved
func (u *User) Solved() revel.Result {
	username := u.Session[USERNAME]
	user := models.GetCurrentUser(username)
	if user != nil {
		//user solve problems
		usps, err := models.FindSovledProblems(user.Id)
		if err != nil {
			u.Flash.Error(err.Error())
			return u.Render()
		} else {
			return u.Render(usps)
		}
	} else {
		return u.Render()
	}

}
