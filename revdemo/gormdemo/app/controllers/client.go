package controllers

import (
	"github.com/revel/revel"
)

type Client struct {
	ResultController
}

func (c Client) Find() revel.Result {
	return c.R.Ok(map[string]interface{}{
		"info": "info",
	})
}
