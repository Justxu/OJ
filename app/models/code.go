package models

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

/*
new error status can noly be added
from the bottom because the status number
in database can not be changed
*/
const (
	Accept int = iota
	CompileError
	WrongAnswer
	TimeLimitExceeded
	MemoryLimitExceeded
	UnHandled
	Handling
	RuntimeError
	PresentationError
	PanicError
)
const (
	Go int = iota
	C
	CPP

	DELIM = "!-_-\n" //delimiter of tests
)

var (
	StatusMap = map[int]string{
		Accept:              "Accept",
		CompileError:        "Compile Error",
		WrongAnswer:         "Wrong Answer",
		TimeLimitExceeded:   "Time Limit Exceeded",
		MemoryLimitExceeded: "Memory Limit Exceeded",
		UnHandled:           "Unhandled",
		Handling:            "Handling",
		RuntimeError:        "Runtime Error",
		PresentationError:   "Present Error",
		PanicError:          "Panic Error",
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
	Id          int64
	UserId      int64
	ProblemId   int64
	CreatedAt   time.Time
	Time        time.Duration
	Status      int
	Lang        int    //source file language
	Memory      int64  //Kb
	Path        string //file path
	Nth         int    //the number of the test not passed
	WrongAnswer string //the last wrong answer
	PanicError  string //panic ouput

}

//check report
type Report struct {
	Tests []Test //all tests
	Nth   int    //nth test is wrong, if Nth is 0 , all passed
}
type Test struct {
	In  string
	Out string
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
func check(err error) {
	if err != nil {
		panic(err)
	}
}
func (s *Source) View() (string, error) {
	source := new(Source)
	_, err := engine.Id(s.Id).Get(source)
	if err != nil {
		return "", err
	}
	in, err := os.Open(path.Join(s.Path, "tmp."+LangMap[s.Lang]))
	check(err)
	defer in.Close()
	code, err := ioutil.ReadAll(in)
	check(err)
	return fmt.Sprintf("%s", code), nil
}

func (s *Source) Check() (*Report, error) {
	source := new(Source)
	_, err := engine.Id(s.Id).Get(source)
	if err != nil {
		return nil, err
	}
	p := new(Problem)
	_, err = engine.Id(source.ProblemId).Get(p)
	if err != nil {
		return nil, err
	}
	in, err := os.Open(p.InputTestPath)
	check(err)
	defer in.Close()
	out, err := os.Open(p.OutputTestPath)
	check(err)
	defer out.Close()
	input, err := ioutil.ReadAll(in)
	check(err)
	output, err := ioutil.ReadAll(out)
	check(err)
	inputs := bytes.Split(input, []byte(DELIM))
	outputs := bytes.Split(output, []byte(DELIM))
	var report Report
	fmt.Printf("%s\n%s\n\n\n%d\n", input, output, source.Nth)
	for i := 0; i < source.Nth; i++ {
		in := fmt.Sprintf("%s", inputs[i])
		out := fmt.Sprintf("%s", outputs[i])
		report.Tests = append(report.Tests, Test{in, out})
	}
	if source.Nth != 0 {
		report.Tests[len(report.Tests)-1].Out = source.WrongAnswer
	}
	report.Nth = source.Nth
	return &report, nil
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
