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
	host, _ := c.String("postgres", "host")
	port, _ := c.String("postgres", "port")
	dataSource := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s", user, password, dbname, sslmode, host, port)
	engine, err = xorm.NewEngine("postgres", dataSource)
	if err != nil {
		panic(err)
	}
	showSQL, _ := c.Bool("postgres", "show_sql")
	/*
		showErr, _ := c.Bool("postgres", "show_err")
		showDebug, _ := c.Bool("postgres", "show_debug")
		showWarn, _ := c.Bool("postgres", "show_warn")
	*/
	engine.ShowSQL = showSQL
	fmt.Println("show SQL", showSQL)
	/*
		engine.ShowWarn = showWarn
		engine.ShowErr = showErr
		engine.ShowDebug = showDebug
	*/
	err = engine.Sync(
		new(Source),
		new(Problem),
		new(User),
		new(Solve),
	)
	if err != nil {
		panic(err)
	}
}
