package main

import (
	"github.com/chaosannals/trial-go/controllers"
	"github.com/chaosannals/trial-go/logics"
	"github.com/gin-gonic/gin"
)

func main() {
	logics.Init()

	defer logics.Recover()

	r := gin.Default()
	
	r.PUT("/insert", controllers.Insert)
	r.GET("/search", controllers.Search)
	r.POST("/update", controllers.Update)
	r.DELETE("/remove", controllers.Remove)
	r.POST("/reword", controllers.Reword)

	r.Run()
}
