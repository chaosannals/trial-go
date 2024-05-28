package main

import (
	"github.com/chaosannals/fulltext/leveldbdemo2/httpapi"
	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
)

func init() {
	keydb.InitSeg()
}

func main() {
	deinitDb := keydb.InitDb()
	defer deinitDb()

	r := gin.Default()
	r.GET("/ping", httpapi.Ping)
	r.PUT("/add", httpapi.Add)
	r.PUT("/add_batch", httpapi.AddBatch)
	r.GET("/query", httpapi.Query)

	r.Run("127.0.0.1:23456")
}
