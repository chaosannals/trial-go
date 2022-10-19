package main

import (
	"log"

	"github.com/chaosannals/trial-go-digdemo/trial"
	"go.uber.org/dig"
)

func main() {
	c := dig.New()
	if err := c.Provide(trial.NewConf); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(trial.NewTcpServer); err != nil {
		log.Fatal(err)
	}
	if err := c.Provide(trial.NewZLogger); err != nil {
		log.Fatal(err)
	}
	if err := c.Invoke(func(server *trial.TcpServer) {
		server.Serve()
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("end.")
}
