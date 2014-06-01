package models

import (
	"fmt"

	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
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
	c, err := config.ReadDefault("conf/misc.conf")
	if err != nil {
		panic(err)
	}
	if c == nil {
		panic("conf path not founded.")
	}
	user, _ := c.String("postgres", "user")
	password, _ := c.String("postgres", "password")
	dbname, _ := c.String("postgres", "dbname")
	sslmode, _ := c.String("postgres", "sslmode")
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)
	engine, err = xorm.NewEngine("postgres", dataSource)
	if err != nil {
		panic(err)
	}
	showSQL, _ := c.Bool("postgres", "show_sql")
	engine.ShowSQL = showSQL
	err = engine.Sync(
		new(Source),
		new(Problem),
		new(User))
	if err != nil {
		panic(err)
	}
}
