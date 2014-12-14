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
	c := judge.Judge("c", "tests/test/A1/src", "tests/test/A1/src/input", "tests/test/A1/src/output", 1000, 10000)

	t.AssertEqual(models.Accept, c.Status)
	c = judge.Judge("go", "tests/test/A1/src", "tests/test/A1/src/input", "tests/test/A1/src/output", 1000, 10000)
	t.AssertEqual(models.Accept, c.Status)
}
