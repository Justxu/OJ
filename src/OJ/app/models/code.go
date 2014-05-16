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
	UnHandled
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
	CreatedAt time.Time
	Status    int
	Path      string
}

func (s *Source) GenPath() string {
	s.Path = "code/" + UUPath()
	return s.Path
}
func (s *Source) StatusString() string {
	return StatusMap[s.Status]
}
