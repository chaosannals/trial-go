package httpapi

import (
	"fmt"
	"net/http"

	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
)

func Add(ctx *gin.Context) {
	doc := &keydb.DocContent{}
	if err := ctx.Bind(doc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	fmt.Printf("add: %v\n", doc.Content)

	if r, err := doc.InsertAndCut(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"result":  r,
			"message": "ok",
		})
	}
}

func AddBatch(ctx *gin.Context) {
	docs := []keydb.DocContent{}
	if err := ctx.Bind(&docs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if result, err := keydb.AddBatch(docs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "ok",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"result":  result,
			"message": "ok",
		})
	}
}
