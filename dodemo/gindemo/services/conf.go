package services

import "github.com/samber/do"

type ConfService struct {
	Host string
	Port int32
}

func NewConfService(i *do.Injector) (*ConfService, error) {
	// 读文件
	return &ConfService{
		Host: "0.0.0.0",
		Port: 44444,
	}, nil
}
