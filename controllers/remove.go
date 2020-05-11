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
	var e = c.BindJSON(&param)
	if e == nil {
		models.Remove(param.Id)
		c.JSON(200, gin.H{
			"message": "ok",
		})
	} else {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
}
