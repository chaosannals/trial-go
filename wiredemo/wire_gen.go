// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/chaosannals/trial-go-wiredemo/trial"
)

// Injectors from wire.go:

func InitTrial() (*trial.TcpServer, error) {
	conf, err := trial.NewConf()
	if err != nil {
		return nil, err
	}
	tcpServer, err := trial.NewTcpServer(conf)
	if err != nil {
		return nil, err
	}
	return tcpServer, nil
}
