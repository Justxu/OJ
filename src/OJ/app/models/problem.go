package models

import (
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

type Problem struct {
	Id          int64
	Title       string
	Solved      int64
	Description string `xorm:"TEXT"` //问题描述
	InputTest   string //输入测试
	OutputTest  string //输出测试
}

func (p *Problem) TestPath() string {
	ui := uuid.NewUUID()
	path := strings.Replace(ui.String(), "-", "", -1)
	return "problem/" + path
}
