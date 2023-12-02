package services

import "github.com/samber/do"

type ConfService struct {
	HttpHost string
	HttpPort uint16
	DbHost   string
	DbPort   uint16
	DbName   string
	DbUser   string
	DbPass   string
}

func NewConfService(i *do.Injector) (*ConfService, error) {
	return &ConfService{
		HttpHost: "127.0.0.1",
		HttpPort: 44400,
		DbHost:   "localhost",
		DbPort:   3306,
		DbName:   "demo",
		DbUser:   "root",
		DbPass:   "123456",
	}, nil
}
