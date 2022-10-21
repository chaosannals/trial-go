package controllers

import "github.com/revel/revel"

type Visit struct {
	*revel.Controller
}

func (c Visit) Index() revel.Result {
	p := c.Params.Query.Get("p")
	return c.RenderJSON(map[string]interface{}{
		"aaa": 123,
		"p":   p,
	})
}

func (c Visit) Info() revel.Result {
	vid := c.Params.Route.Get("id")
	return c.RenderJSON(map[string]interface{}{
		"aaa": 123,
		"id":  vid,
	})
}

func (c Visit) Add() revel.Result {
	name := c.Params.Form.Get("name")
	return c.RenderJSON(map[string]interface{}{
		"aaa":  123,
		"name": name,
	})
}

func (c Visit) Edit() revel.Result {
	var data map[string]interface{}
	c.Params.BindJSON(&data)
	return c.RenderJSON(map[string]interface{}{
		"aaa":  123,
		"data": data,
	})
}

func (c Visit) Upload() revel.Result {
	f := c.Params.Files["file_name"]
	return c.RenderJSON(map[string]interface{}{
		"aaa":  123,
		"name": f[0].Filename,
	})
}

// 自动通过参数匹配，类型转换，没有为该类型的 zero value
func (c Visit) Delete(id int32, number string) revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"aaa":    123,
		"id":     id,
		"number": number,
	})
}
