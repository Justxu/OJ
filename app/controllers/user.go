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
	err := engine.Limit(10).Find(&users)
	if err != nil {
		fmt.Println(err)
	}
	return u.Render(users)
}

func (u *User) Profile() revel.Result {
	var user models.User
	username := u.Session["username"]
	engine.Where("name = ?", username).Get(&user)
	return u.Render(user)
}
