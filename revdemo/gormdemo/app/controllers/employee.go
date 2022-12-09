package controllers

import (
	"gormdemo/app"
	"gormdemo/entities"
	"gormdemo/models"
	"gormdemo/util"
	"math"
	"net/http"

	"github.com/revel/revel"
)

type Employee struct {
	*revel.Controller
}

type AA struct {
	I1 int `json:"i1"`
	I2 int `json:"i2"`
	F2 float32 `json:"f2"`
}

func (c Employee) Find(p int, ps int) revel.Result {
	// gen 的 models 是常规查询方式
	psize := (int)(math.Min(100, (float64)(ps)))
	offset := (int)(math.Max((float64)(p - 1), 0)) * psize
	rows, err := app.Db.
		Model(&models.EEmployee{}).
		Where("removed_at IS NULL").
		Offset(offset).
		Limit(psize).
		Rows()

	util.De[AA]("asdf")

	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]interface{}{
			"err": err,
		})
	}

	a, err := util.JsonAs[AA]("{\"i1\": 123, \"i2\": 123, \"f2\": 123.12}")

	return c.RenderJSON(map[string]interface{}{
		"rows": rows,
		"a": a,
		"err": err,
	})
}

func (c Employee) Info(id uint64) revel.Result {
	var err error
	// gorm gen 生成的 entities 是另一种查询方式。
	info, err := entities.EEmployee.
		WithContext(c.Request.Context()).
		Where(entities.EEmployee.ID.Eq(id)).
		First()
	if err != nil {
		return c.RenderError(err)
	}
	return c.RenderJSON(map[string]interface{}{
		"info": info,
	})
}
