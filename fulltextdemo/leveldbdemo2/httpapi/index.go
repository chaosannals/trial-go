package httpapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

type NgramParam struct {
	Plain string `json:"plain" form:"plain"`
	Min   uint8  `json:"min" form:"min"`
	Max   uint8  `json:"max" form:"max"`
}

func Ngram(ctx *gin.Context) {
	param := &NgramParam{}
	if err := ctx.BindJSON(param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}
	pl := uint8(len(param.Plain))
	seqs := make([]string, 0)
	for i := param.Min; i <= param.Max; i++ {
		end := pl - i
		for j := uint8(0); j < end; j++ {
			seg := param.Plain[j : j+i]
			seqs = append(seqs, seg)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"seqs": seqs,
	})
}
