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
func (controller *IndexController) AddAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) DeleteAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) UpdateAccount(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) UploadAccountImage(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{})
}
