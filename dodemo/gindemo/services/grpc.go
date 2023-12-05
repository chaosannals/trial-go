package services

import (
	"context"
	"log/slog"

	"github.com/chaosannals/dodemo-gindemo/apis"
	"github.com/samber/do"
	"google.golang.org/grpc"
)

func NewGrpcServer(i *do.Injector) (*grpc.Server, error) {
	server := grpc.NewServer()
	logger := do.MustInvoke[*slog.Logger](i)
	apis.RegisterHelloServiceServer(server, &ApisHelloServer{
		logger: logger,
	})
	return server, nil
}

type ApisHelloServer struct {
	apis.HelloServiceServer
	logger *slog.Logger
}

func (i *ApisHelloServer) Hello(ctx context.Context, request *apis.HelloRequest) (*apis.HelloResponse, error) {
	i.logger.Debug("call hello")

	return &apis.HelloResponse{
		Message: request.Message,
	}, nil
}
