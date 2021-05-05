package controllers

import (
	"github.com/chaosannals/trial-go/logics"
	"github.com/chaosannals/trial-go/models"
	"github.com/gin-gonic/gin"
)

//InsertRequestParam 插入参数
type InsertRequestParam struct {
	ID      string `json:"id"`      //ID
	TypeID  uint64 `json:"type_id"` // 类型ID
	Title   string `json:"title"`   // 标题
	Content string `json:"content"` // 内容
}

//Insert 插入文档
func Insert(c *gin.Context) {
	var param InsertRequestParam
	if e := c.BindJSON(&param); e != nil {
		c.JSON(400, gin.H{
			"message": e.Error(),
		})
	}
	logics.Insert(models.Store{
		ID:      param.ID,
		TypeID:  param.TypeID,
		Title:   param.Title,
		Content: param.Content,
	})
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
