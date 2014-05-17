package models

type Problem struct {
	Id     int64 `xorm:"pk"`
	Name   string
	Solved int
	//问题描述
	Des string `xorm:"TEXT"`
}
