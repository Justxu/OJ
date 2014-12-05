package tests

import "github.com/revel/revel"
import "regexp"

type OutputTest struct {
	revel.TestSuite
}

type testPair struct {
	result bool
	input  string
}

func (t *OutputTest) TestLegalOutput() {
	var rep = `\w\w:\d+:\d+:([\s\S]*)`

	var tests = []testPair{{true, `WA:0:0:asd
    asdasd
    qweqwewq   
    qweqwe\e1231231`}, {true, `LE:0:0:asd
    asd`}, {false, `LMT:0:0asd`},
	}
	for _, v := range tests {
		t.Assert(regexp.MustCompile(rep).MatchString(v.input) == v.result)
	}
}
