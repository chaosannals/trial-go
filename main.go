package main

import (
	"github.com/gin-gonic/gin"
	"github.com/chaosannals/trial-go/controllers"
	"github.com/chaosannals/trial-go/logics"
)

func main() {
	rd := logics.Init()

	defer rd()

	r := gin.Default()
	r.PUT("/change", controllers.Change)
	r.GET("/search", controllers.Search)
	r.DELETE("/remove", controllers.Remove)

	r.Run()
}
