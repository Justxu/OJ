package models

type Problem struct {
	Id          int64
	Title       string
	Solved      int64
	Description string `xorm:"TEXT"` //问题描述
	InputTest   string //输入测试
	OutputTest  string //输出测试
}
