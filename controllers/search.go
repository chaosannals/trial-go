package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-ego/riot/types"
	"github.com/chaosannals/trial-go/models"
)

func Search(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": models.Search(types.SearchReq{Text:"复仇者"}),
	})
}