package controllers

import (
	"database/sql"
	"fmt"
	"os"
	"reveldemo/app"
	"reveldemo/app/models"
	"strconv"

	"github.com/revel/revel"
)

type Visit struct {
	*revel.Controller
}

func (c Visit) Index() revel.Result {
	p := c.Params.Query.Get("p")
	s := c.Params.Query.Get("s")
	var pi int64
	var si int64
	var err error
	if pi, err = strconv.ParseInt(p, 10, 64); err != nil {
		pi = 1
	}
	if si, err = strconv.ParseInt(s, 10, 64); err != nil {
		si = 20
	}
	if si > 100 {
		si = 100
	}

	q := fmt.Sprintf("SELECT * FROM rd_visit LIMIT %d, %d", (pi-1)*si, si)
	var rows *sql.Rows
	if rows, err = app.DB.Query(q); err != nil {
		return c.RenderJSON(map[string]interface{}{
			"p":   p,
			"s":   s,
			"err": err,
		})
	}
	return c.RenderJSON(map[string]interface{}{
		"p":    p,
		"s":    s,
		"rows": rows,
	})
}

func (c Visit) Info() revel.Result {
	vid := c.Params.Route.Get("id")
	var vidi int64
	var err error
	if vidi, err = strconv.ParseInt(vid, 10, 64); err != nil {
		return c.RenderJSON(map[string]interface{}{
			"id":  vid,
			"msg": "ID 非数字",
		})
	}

	var info models.VisitModel
	if err = app.DBM.SelectOne(&info, "SELECT * FROM rd_visit WHERE id=?", vidi); err != nil {
		return c.RenderJSON(map[string]interface{}{
			"id":  vid,
			"msg": "不是有效的ID",
		})
	}
	return c.RenderJSON(map[string]interface{}{
		"id":   vid,
		"info": &info,
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

func (c Visit) Download() revel.Result {
	f, err := os.Open("README.md")
	if err != nil {
		return c.RenderFileName("README.md", revel.Attachment)
	}
	return c.RenderFile(f, revel.Inline)
}

// 自动通过参数匹配，类型转换，没有为该类型的 zero value
func (c Visit) Delete(id int32, number string) revel.Result {
	return c.RenderJSON(map[string]interface{}{
		"aaa":    123,
		"id":     id,
		"number": number,
	})
}
