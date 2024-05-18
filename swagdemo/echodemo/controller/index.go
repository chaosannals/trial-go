package controller

import (
	"errors"
	"net/http"

	"github.com/chaosannals/swagdemo/echodemo/util"
	"github.com/labstack/echo/v4"
)

type IndexController struct {
}

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  any
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /accounts/{id} [get]
func (controller *IndexController) ShowAccount(ctx echo.Context) error {
	util.NewError(ctx, http.StatusUnauthorized, errors.New("Authorization is required Header"))
	return ctx.JSON(http.StatusOK, map[string]any{})
}

// ListAccounts godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  any
// @Failure      400  {object}  util.HTTPError
// @Failure      404  {object}  util.HTTPError
// @Failure      500  {object}  util.HTTPError
// @Router       /accounts/list [get]
func (controller *IndexController) ListAccounts(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}

type IndexAddParam struct {
	Page     int `json:"page" minimum:"1" validate:"optional" example:"1"`
	PageSize int `json:"page_size" minimum:"1" maximum:"20" validate:"optional" example:"10"`
}

// 特殊名  request 被使用时，指定整个参数体结构 @Param

// @Summary     Some Endpoint
// @Produce     json
// @Param       request             query    IndexAddParam false "Query Params"
// @Router       /accounts/add [put]
func (controller *IndexController) AddAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}

// @Summary     Some Endpoint
// @Produce     json
// @Param       request             body    IndexAddParam false "Query Params"
// @Router       /accounts/delete [delete]
func (controller *IndexController) DeleteAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) UpdateAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) UploadAccountImage(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
