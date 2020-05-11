package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ego/riot/types"
	"github.com/chaosannals/trial-go/models"
)

func Change(c *gin.Context) {
	text := "《复仇者联盟3：无限战争》是全片使用IMAX摄影机拍摄"
	text1 := "在IMAX影院放映时"
	text2 := "全片以上下扩展至IMAX 1.9：1的宽高比来呈现"
	// 将文档加入索引，docId 从1开始
	models.Change("1", types.DocData{Content: text})
	models.Change("2", types.DocData{Content: text1}, false)
	models.Change("3", types.DocData{Content: text2}, true)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
