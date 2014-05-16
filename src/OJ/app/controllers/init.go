package controllers

import (
	"fmt"

	"github.com/go-xorm/xorm"

	"OJ/app/models"
)

var (
	engine *xorm.Engine
)

func init() {
	fmt.Println("init")
	engine = models.Engine()
}
