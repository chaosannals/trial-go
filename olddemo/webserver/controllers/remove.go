package controllers

import (
	"github.com/chaosannals/trial-go/logics"
	"github.com/gin-gonic/gin"
)

//RemoveRequestParam 删除参数
type RemoveRequestParam struct {
	ID string `json:"id"` // ID
}

//Remove 删除
func Remove(c *gin.Context) {
	var param RemoveRequestParam
	if e := c.BindJSON(&param); e != nil {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	logics.Remove(param.ID)
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
