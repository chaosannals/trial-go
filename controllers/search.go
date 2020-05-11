package controllers

import (
	"github.com/chaosannals/trial-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-ego/riot/types"
)

type SearchRequestParam struct {
	Text string `json:"text"`
}

func Search(c *gin.Context) {
	var param SearchRequestParam
	var e = c.BindJSON(&param)
	if e == nil {
		c.JSON(200, gin.H{
			"message": models.Search(types.SearchReq{Text: param.Text}),
		})
	} else {
		c.JSON(200, gin.H{
			"message": e.Error(),
		})
	}
}
