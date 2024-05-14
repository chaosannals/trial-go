package controller

import (
	"errors"
	"net/http"

	"github.com/chaosannals/swagdemo/gindemo/util"
	"github.com/gin-gonic/gin"
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
func (controller *IndexController) ShowAccount(ctx *gin.Context) {
	util.NewError(ctx, http.StatusUnauthorized, errors.New("Authorization is required Header"))
	ctx.JSON(http.StatusOK, map[string]any{})
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
func (controller *IndexController) ListAccounts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) AddAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) DeleteAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) UpdateAccount(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]any{})
}
func (controller *IndexController) UploadAccountImage(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]any{})
}
