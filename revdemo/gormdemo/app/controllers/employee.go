package controllers

import (
	"gormdemo/app"
	"gormdemo/models"
	"net/http"

	"github.com/revel/revel"
)

type Employee struct {
	*revel.Controller
}

func (c Employee) Find() revel.Result {
	rows, err := app.Db.Model(&models.EEmployee{}).Rows()
	if err != nil {
		c.Response.Status = http.StatusBadRequest
		return c.RenderJSON(map[string]interface{}{
			"err": err,
		})
	}
	return c.RenderJSON(map[string]interface{}{
		"rows": rows,
	})
}