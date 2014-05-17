package models

import (
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

const (
	Accept int = iota
	CompileError
	WrongAnswer
	TimeLimitExceeded
	MemoryLimitExceeded
	UnHandled
)
const (
	Go int = iota
	C
	CPlusPlus
)

var StatusMap = map[int]string{
	Accept:       "Accecpt",
	CompileError: "CompileError",
}

func UUPath() string {
	ui := uuid.NewUUID()
	path := strings.Replace(ui.String(), "-", "", -1)
	return path + "/main.go"
}

type Source struct {
	Id        int64
	UserId    int64
	ProblemId int64
	CreatedAt time.Time
	Lang      int64
	Status    int
	Time      time.Duration
	//kb为单位
	Memory int64
	//文件路径
	Path string
	//成功测试输入的第n行
	TestLine int
}

func (s *Source) GenPath() string {
	s.Path = "code/" + UUPath()
	return s.Path
}
func (s *Source) StatusString() string {
	return StatusMap[s.Status]
}
