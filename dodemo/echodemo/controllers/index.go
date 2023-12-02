package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type IndexController struct {

}

func NewIndexController(i *do.Injector) (*IndexController, error) {
	return &IndexController{}, nil
}

func (i *IndexController) Index(c echo.Context) error{
	return c.JSON(
		http.StatusOK,
		map[string]any {
			"code": 0,
			"message": "ok",
		},
	)
}