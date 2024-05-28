package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type DocContent struct {
	Content string `json:"content" form:"content"`
}

func (doc *DocContent) Encode() ([]byte, error) {
	return json.Marshal(doc)
}

func Add(ctx *gin.Context) {
	doc := DocContent{}
	if err := ctx.Bind(&doc); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	id, err := uuid.NewUUID()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	de, err := doc.Encode()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if err := keydb.LDB.Put(id[:], de, &opt.WriteOptions{}); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":      id,
		"message": "ok",
	})
}

func AddBatch(ctx *gin.Context) {
}
