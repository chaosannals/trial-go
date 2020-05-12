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
	if e := c.BindJSON(&param); e != nil {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	c.JSON(200, gin.H{
		"message": models.Search(types.SearchReq{Text: param.Text}),
	})
}
