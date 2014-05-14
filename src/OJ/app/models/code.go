package models

import (
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/go-xorm/xorm"
)

func UUPath() string {
	ui := uuid.NewUUID()
	ui = strings.Replace(ui, "-", "", -1)
	return ui + "/main.go"
}

type Source struct {
	Id        int64
	UserId    int64
	CreatedAt time.Time
	Path      string
}

func (s *Source) GenPath() string {
	s.Path = "code/" + UUPath()
	return s.Path
}
