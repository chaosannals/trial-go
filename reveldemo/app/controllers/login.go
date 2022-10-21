package controllers

import (
	"github.com/revel/revel"
)

type Login struct {
	*revel.Controller
}

func (c Login) SignIn() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"msg": "signin",
	})
}

func (c Login) SignOut() revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"msg": "signout",
	})
}

func (c Login) Index() revel.Result {
	return c.Render()
}
