package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/chaosannals/trial-go/logics"
	"github.com/chaosannals/trial-go/bases"
)

func main() {
	logics.Init()

	defer logics.Recover()

	server := bases.NewHttpServer()

	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
	}
}
