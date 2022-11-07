package controllers

import (
	"net/http"

	"github.com/chaosannals/trial-go-digdemo/models"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type EmployeeController struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewEmployeeController(
	db *gorm.DB,
	logger *zerolog.Logger,
) *EmployeeController {
	return &EmployeeController{
		db:     db,
		logger: logger,
	}
}

func (i *EmployeeController) List(c echo.Context) (err error) {
	var rows []models.EEmployee

	i.db.Model(&models.EEmployee{}).Select("*").Find(&rows)

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"rows":    rows,
			"message": "ok",
		},
	)
}

func (i *EmployeeController) Add(c echo.Context) (err error) {
	row := &models.EEmployee{}
	if err = c.Bind(row); err != nil {
		return
	}

	i.db.Create(row)
	err = i.db.Error
	if err != nil {
		return
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"id":      row.ID,
			"message": "ok",
		},
	)
}

func (i *EmployeeController) Edit(c echo.Context) (err error) {
	row := &models.EEmployee{}
	if err = c.Bind(row); err != nil {
		return
	}

	i.db.Model(&models.EEmployee{}).
		Where("id = ?", row.ID).
		Updates(row)
	err = i.db.Error
	if err != nil {
		return
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"id":      row.ID,
			"message": "ok",
		},
	)
}

func (i *EmployeeController) Del(c echo.Context) (err error) {
	id := c.Param("id")
	i.db.Delete(&models.EEmployee{}, id)
	err = i.db.Error
	if err != nil {
		return
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"code":    0,
			"message": "ok",
		},
	)
}
