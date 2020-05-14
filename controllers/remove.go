package controllers

import (
	"github.com/chaosannals/trial-go/logics"
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
	logics.Remove(param.Id)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
