package models

import (
	"fmt"
	"path/filepath"

	"github.com/go-xorm/xorm"
	"github.com/revel/config"
)

var (
	engine *xorm.Engine
)

func Engine() *xorm.Engine {
	return engine
}
func init() {
	var err error
	path, _ := filepath.Abs("")
	c, err := config.ReadDefault(fmt.Sprintf("%s/src/finances/conf/misc.conf", path))
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic("conf path not founded.")
	}
	user, _ := c.String("postgres", "user")
	password, _ := c.String("postgres", "user")
	dbname, _ := c.String("postgres", "dbname")
	sslmode, _ := c.String("postgres", "sslmode")
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	engine, err = xorm.NewEngine("postgres", dataSource)
	if err != nil {
		panic(err)
	}
	showSQL, _ := c.Bool("postgres", "show_sql")
	err = engine.Sync(new(Code))
	if err != nil {
		panic(err)
	}
}
