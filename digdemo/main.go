package main

import (
	"log"
	"net/http"

	"github.com/chaosannals/trial-go-digdemo/trial"
	"github.com/rs/zerolog"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()
	if err := c.Provide(trial.NewConf); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(trial.NewGormDb); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(trial.NewTcpServer); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(trial.NewZLogger); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(trial.NewEchoHttpServer); err != nil {
		log.Fatal(err)
	}

	if err := c.Invoke(func(server *http.Server, logger *zerolog.Logger) {
		if err := server.ListenAndServe(); err != nil {
			logger.Err(err).Msg("http server serve error")
		}
	}); err != nil {
		log.Fatal(err)
	}

	if err := c.Invoke(func(server *trial.TcpServer) {
		server.Serve()
	}); err != nil {
		log.Fatal(err)
	}

	log.Println("end.")
}
