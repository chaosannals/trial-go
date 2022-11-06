package bases

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/chaosannals/trial-go/controllers"
)

func NewHttpServer() *http.Server {
	router := gin.Default()
	router.PUT("/insert", controllers.Insert)
	router.GET("/search", controllers.Search)
	router.POST("/update", controllers.Update)
	router.DELETE("/remove", controllers.Remove)
	router.POST("/reword", controllers.Reword)
	return &http.Server {
		Addr: "0.0.0.0:8080",
		Handler: router,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1000000,
	}
}