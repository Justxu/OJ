package models

type Problem struct {
	Id          int64
	Title       string
	Solved      int64
	Description string `xorm:"TEXT"` //问题描述
}
