package controllers

import (
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
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	c.JSON(200, gin.H{
		"message": logics.Search(types.SearchReq{Text: param.Text}),
	})
}
