package controllers

import (
	"github.com/revel/revel"
)

type Sign struct {
	*revel.Controller
}

func (c Sign) Login() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"msg": "login",
	})
}
