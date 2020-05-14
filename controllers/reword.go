package controllers

import (
	"github.com/gin-gonic/gin"
)

//RewordRequestParam 重置索引 
type RewordRequestParam struct {
	ID string `json:id`
}

//Reword 重置文档索引
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
