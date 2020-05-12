package controllers

import (
	"github.com/chaosannals/trial-go/models"
	"github.com/gin-gonic/gin"
)

type RemoveRequestParam struct {
	Id string `json:"id"`
}

func Remove(c *gin.Context) {
	var param RemoveRequestParam
	if e := c.BindJSON(&param); e != nil {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	models.Remove(param.Id)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
