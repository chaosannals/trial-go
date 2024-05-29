package httpapi

import (
	"fmt"
	"net/http"

	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	id, seqs, err := doc.InsertAndCut()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":      id,
		"seqs":    seqs,
		"message": "ok",
	})
}

func AddBatch(ctx *gin.Context) {
	docs := []keydb.DocContent{}
	if err := ctx.Bind(&docs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	result := make(map[uuid.UUID][]string)

	for _, doc := range docs {
		id, seqs, err := doc.InsertAndCut()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
			return
		}
		result[id] = seqs
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result":  result,
		"message": "ok",
	})
}
