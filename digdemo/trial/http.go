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
	apiEmployeeGroup.GET("/list", controllers.EmployeeList)
	apiEmployeeGroup.PUT("/add", controllers.EmployeeAdd)
	apiEmployeeGroup.DELETE("/delete", controllers.EmployeeDel)
	//apiGroup.GET("")

	for _, r := range e.Routes() {
		logger.Info().
			Str("name", r.Name).
			Str("path", r.Path).
			Str("method", r.Method).
			Msg("route:")
	}

	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", conf.HttpHost, conf.HttpPort),
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1000000,
	}
	return server
}
