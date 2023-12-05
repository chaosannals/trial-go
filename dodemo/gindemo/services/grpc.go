package services

import (
	"context"

	"github.com/chaosannals/dodemo-gindemo/apis"
	"github.com/samber/do"
	"google.golang.org/grpc"
)

func NewGrpcServer(i *do.Injector) (*grpc.Server, error) {
	server := grpc.NewServer()
	apis.RegisterHelloServiceServer(server, &ApisHelloServer{})
	return server, nil
}

type ApisHelloServer struct {
	apis.HelloServiceServer
}

func (i *ApisHelloServer) Hello(ctx context.Context, request *apis.HelloRequest) (*apis.HelloResponse, error) {
	return &apis.HelloResponse{
		Message: request.Message,
	}, nil
}
