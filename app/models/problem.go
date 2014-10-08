package models

import (
	"strings"

	"code.google.com/p/go-uuid/uuid"
)

type Problem struct {
	Id             int64
	Title          string
	Solved         int64 //times of accepted submit
	TimeLimit      int64
	MemoryLimit    int64
	Description    string `xorm:"TEXT"`
	InputSample    string `xorm:"varchar(512)"`
	OutputSample   string `xorm:"varchar(512)"`
	InputTestPath  string //input test path
	OutputTestPath string //output test path
	ImgSrc         string //url
}

func (p *Problem) TestPath() string {
	ui := uuid.NewUUID()
	path := strings.Replace(ui.String(), "-", "", -1)
	return "problem/" + path
}
