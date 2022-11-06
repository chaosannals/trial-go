package controllers

import (
	"net/http"

	"github.com/chaosannals/trial-go/logics"
	"github.com/gin-gonic/gin"
	"github.com/go-ego/riot/types"
)

//SearchRequestParam 搜索参数
type SearchRequestParam struct {
	Text string `json:"text"`
}

//Search 搜索
func Search(c *gin.Context) {
	var param SearchRequestParam
	if e := c.BindJSON(&param); e != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": e.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": logics.Search(types.SearchReq{Text: param.Text}),
	})
}
