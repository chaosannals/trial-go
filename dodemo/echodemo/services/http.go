package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaosannals/dodemo-echodemo/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/samber/do"
)

type EchoHttpService struct {
	server *http.Server
	logger *zerolog.Logger
}

func NewEchoHttpService(i *do.Injector) (*EchoHttpService, error) {
	conf := do.MustInvoke[*ConfService](i)
	logger := do.MustInvoke[*zerolog.Logger](i)
	indexController := do.MustInvoke[*controllers.IndexController](i)

	e := echo.New()
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		logger.Err(err).Int("status code", code).Msg("http error:")
		c.JSON(
			http.StatusBadRequest,
			map[string]any{
				"code":    -1,
				"message": "请求错误",
			},
		)
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", indexController.Index)

	for _, r := range e.Routes() {
		logger.Info().
			Str("name", r.Name).
			Str("path", r.Path).
			Str("method", r.Method).
			Msg("route:")
	}

	address := fmt.Sprintf("%s:%d", conf.HttpHost, conf.HttpPort)
	logger.Info().Str("address", address).Msg("http server bind:")
	server := &http.Server{
		Addr:           address,
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1000000,
	}
	return &EchoHttpService{
		server: server,
		logger: logger,
	}, nil
}

func (i *EchoHttpService) HealthCheck() error {
	return fmt.Errorf("engine broken")
}

func (i *EchoHttpService) Serve() {
	fmt.Println("start serve")
	if err := i.server.ListenAndServe(); err != nil {
		i.logger.Err(err).Msg("http server serve error")
	}
}
