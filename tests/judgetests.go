package tests

import (
	"github.com/ggaaooppeenngg/OJ/app/judge"
	"github.com/ggaaooppeenngg/OJ/app/models"

	"github.com/revel/revel"
)

type JudegeTest struct {
	revel.TestSuite
}

func (t *JudegeTest) TestCLangAPlusB() {
	c, e := judge.Judge("c", "tests/test/A1/", "tests/test/A1/src/input", "tests/test/A1/src/output", 1000, 10000)
	t.Assert(e == nil)
	t.Assert(c == models.Accept)
}
