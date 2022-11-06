package trial

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosannals/trial-go-digdemo/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

func NewEchoHttpServer(conf *Conf, logger *zerolog.Logger) *http.Server {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		logger.Err(err).Int("status code", code)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	apiGroup := e.Group("/api")
	apiEmployeeGroup := apiGroup.Group("/employee")
	apiEmployeeGroup.GET("/add", controllers.EmployeeAdd)
	//apiGroup.GET("")

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", conf.HttpHost, conf.HttpPort),
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1000000,
	}
	return server
}
