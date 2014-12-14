package controllers

import (
	"github.com/revel/revel"
)

// ouput errors
type Notice struct {
	*revel.Controller
}

func (c *Notice) Crash() revel.Result {
	return c.Render()
}

//
func (c *Notice) Search() revel.Result {
	return nil
}
