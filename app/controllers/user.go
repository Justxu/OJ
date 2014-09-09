package controllers

import (
	"fmt"

	"github.com/ggaaooppeenngg/OJ/app/models"

	"github.com/revel/revel"
)

type User struct {
	*revel.Controller
}

func (u *User) Rating() revel.Result {
	var users []models.User
	err := engine.Limit(10).Desc("solved").Find(&users)
	if err != nil {
		fmt.Println(err)
	}
	return u.Render(users)
}

func (u *User) ProfileVisit(id int64) revel.Result {
	var user models.User
	engine.Id(id).Get(&user)
	return u.Render(user)
}

func (u *User) Profile() revel.Result {
	var user models.User
	username := u.Session["username"]
	engine.Where("name = ?", username).Get(&user)
	return u.Render(user)
}

func (u *User) Solved() revel.Result {
	username := u.Session["username"]
	user := models.GetCurrentUser(username)
	if user != nil {
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
