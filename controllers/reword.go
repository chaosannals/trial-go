package controllers

import (
	"github.com/gin-gonic/gin"
)

type RewordRequestParam struct {
	Word string `json:word`
}

func Reword(c *gin.Context) {
	var param RemoveRequestParam
	if e := c.BindJSON(&param); e != nil {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
