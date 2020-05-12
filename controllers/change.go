package controllers

import (
	"github.com/chaosannals/trial-go/models"
	"github.com/gin-gonic/gin"
	"github.com/go-ego/riot/types"
)

type ChangeRequestParam struct {
	Id      string `json:"id"`
	Content string `json:"content"`
}

func Change(c *gin.Context) {
	var param ChangeRequestParam
	if e := c.BindJSON(&param); e != nil {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	models.Change(param.Id, types.DocData{Content: param.Content}, true)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
