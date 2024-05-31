package httpapi

import (
	"fmt"
	"net/http"

	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
)

type QueryParam struct {
	Plain string `json:"plain" form:"plain"`
}

func Query(ctx *gin.Context) {
	var param QueryParam
	// gin 框架有 BUG ，同样做法 add.go 的接口可以用 Bind, 这个接口必须 BindJson ...
	// 大概是 GET 请求他不让带 BODY ，但是只要指定 BindJson 时，gin 又允许获取 BODY
	// HTTP 没有规定 GET 不可以带 BODY ，这样违规啊。
	// if err := ctx.Bind(&param); err != nil {
	if err := ctx.BindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	fmt.Printf("Query: %v\n", param)
	result, err := keydb.Query(&keydb.QueryParam{
		Plain: param.Plain,
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
