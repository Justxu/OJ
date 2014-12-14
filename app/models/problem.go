package models

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ggaaooppeenngg/util"

	"code.google.com/p/go-uuid/uuid"
	"github.com/revel/revel"
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
	IsValid        bool   //flag for checking probelm offered by ordinary users,
	PosterId       int64  //Post id
}

func (p *Problem) Validate(v *revel.Validation, in, out []byte) {
	v.Required(p.Title).Message("Title Required")
	v.Min(int(p.MemoryLimit), 1).Message("TimeLimit Required")
	v.Min(int(p.TimeLimit), 1).Message("MemoryLimit Required")
	v.Required(p.Description).Message("Description Required")
	v.Required(in).Message("input file needed")
	v.Required(out).Message("output file needed")
	v.MaxSize(p.InputSample, 512).Message("input sample too long")
	v.MaxSize(p.OutputSample, 512).Message("output sample too long")
	path := p.TestPath()
	p.InputTestPath = path + "/inputTest"
	p.OutputTestPath = path + "/outputTest"
}

func (p *Problem) Delete() error {
	_, err := engine.Id(p.Id).Delete(new(Problem))
	if err != nil {
		return err
	}
	ip := p.InputTestPath
	i := strings.Index(ip, "/inputTest")
	if i > len(ip) {
		log.Println("index out of range")
		log.Println(i, len(ip))
		return errors.New("index out of range")
	} else {
		dir := ip[:i]
		if util.IsExist(dir) {
			return os.RemoveAll(ip[:i])
		} else {
			return nil
		}
	}
}

func (p *Problem) Poster() string {
	user := new(User)
	has, err := engine.Id(p.PosterId).Get(user)
	if err != nil || !has {
		fmt.Println(err)
		return "guest"
	} else {
		return user.Name
	}
}
func (p *Problem) TestPath() string {
	ui := uuid.NewUUID()
	path := strings.Replace(ui.String(), "-", "", -1)
	return "problem/" + path
}
