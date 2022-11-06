package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func EmployeeList(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"message": "ok",
		},
	)
}

func EmployeeAdd(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"message": "ok",
		},
	)
}

func EmployeeDel(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"message": "ok",
		},
	)
}