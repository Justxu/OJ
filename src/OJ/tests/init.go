package tests

import (
	"OJ/app/models"

	"github.com/go-xorm/xorm"
)

var (
	engine *xorm.Engine
)

func init() {
	engine = models.Engine()
}
