package httpapi

import (
	"net/http"

	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
)

func Task(ctx *gin.Context) {
	keydb.DataTaskQueue <- keydb.DataTask{}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
