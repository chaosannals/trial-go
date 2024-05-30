package main

import (
	"log"

	"github.com/chaosannals/fulltext/leveldbdemo2/httpapi"
	"github.com/chaosannals/fulltext/leveldbdemo2/keydb"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	keydb.InitSeg()
	if err := keydb.InitMysql(); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	deinitDb := keydb.InitDb()
	defer deinitDb()

	r := gin.Default()
	r.GET("/ping", httpapi.Ping)
	r.GET("/ngram", httpapi.Ngram)
	r.GET("/query", httpapi.Query)
	r.PUT("/add", httpapi.Add)
	r.PUT("/add_batch", httpapi.AddBatch)
	r.POST("/task", httpapi.Task)

	r.Run("127.0.0.1:23456")
}
