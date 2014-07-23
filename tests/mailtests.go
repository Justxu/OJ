package tests

import (
	"OJ/app/controllers"

	"github.com/revel/revel"
)

type MailTest struct {
	revel.TestSuite
}

func (t MailTest) TestSendMail() {
	stmp := controllers.GetStmp()
	err := controllers.SendMail("subject", "message", "zuoyejizhuce@163.com", []string{"ggaaooppeenngg@qq.com"}, stmp, true)
	t.AssertEqual(err, nil)
}
