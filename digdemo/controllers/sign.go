package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type SignController struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func NewSignController(
	db *gorm.DB,
	logger *zerolog.Logger,
) *SignController {
	return &SignController{
		db:     db,
		logger: logger,
	}
}

func (i *SignController) Login(c echo.Context) error {
	return c.JSON(
		http.StatusOK,
		map[string]any {
			"code": 0,
			"message": "ok",
		},
	)
}
