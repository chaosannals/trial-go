package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chaosannals/dodemo-gindemo/apis"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApisHelloClientService struct {
	logger *slog.Logger
	conf   *ConfService
}

func NewApisHelloClientService(i *do.Injector) (*ApisHelloClientService, error) {
	logger := do.MustInvoke[*slog.Logger](i)
	conf := do.MustInvoke[*ConfService](i)

	return &ApisHelloClientService{
		logger: logger,
		conf:   conf,
	}, nil
}

func (i *ApisHelloClientService) SayHello(msg string) (*apis.HelloResponse, error) {
	address := fmt.Sprintf("127.0.0.1:%d", i.conf.Port)

	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := apis.NewHelloServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return client.Hello(ctx, &apis.HelloRequest{
		Message: msg,
	})
}
