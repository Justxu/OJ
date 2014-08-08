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
	Handling
)
const (
	Go int = iota
	C
	CPP
)

var (
	StatusMap = map[int]string{
		Accept:              "Accept",
		CompileError:        "Compile Error",
		WrongAnswer:         "Wrong Answer",
		TimeLimitExceeded:   "Time Limit Exceeded",
		MemoryLimitExceeded: "Memory Limit Exceeded",
		UnHandled:           "UnHandled",
		Handling:            "Handling",
	}
	LangMap = map[int]string{
		Go:  "go",
		C:   "c",
		CPP: "cpp",
	}
)

func UUPath() string {
	ui := uuid.NewUUID()
	path := strings.Replace(ui.String(), "-", "", -1)
	return path
}

type Source struct {
	Id        int64
	UserId    int64
	ProblemId int64
	CreatedAt time.Time
	Lang      int
	Status    int
	Time      time.Duration
	Memory    int64  //以Kb为单位
	Path      string //文件路劲
	TestLine  int    //测试输入的第N行
}

func (s *Source) TimeUsed() int64 {
	return s.Time.Nanoseconds() / 1000
}
func (s *Source) CreatedTime() string {
	return s.CreatedAt.Format(time.Kitchen)
}
func (s *Source) GetUserName() string {
	u := new(User)
	_, err := engine.Where("id = ?", s.UserId).Cols("name").Get(u)
	if err != nil {
		return err.Error()
	}
	return u.Name

}
func (s *Source) GetProblemTitle() string {
	p := new(Problem)
	_, err := engine.Where("id = ?", s.ProblemId).Cols("title").Get(p)
	if err != nil {
		return err.Error()
	}
	return p.Title
}
func (s *Source) GenPath() string {
	s.Path = "code/" + UUPath()
	return s.Path
}
func (s *Source) StatusString() string {
	return StatusMap[s.Status]
}
func (s *Source) LangString() string {
	return LangMap[s.Lang]
}
