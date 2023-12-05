//go:generate protoc -I=./protos  --go_out=./  --go-grpc_out=./ ./protos/*.proto

package main

import (
	"github.com/chaosannals/dodemo-gindemo/services"
	"github.com/samber/do"
)

func main() {
	// 容器
	injector := do.New()
	defer injector.Shutdown()

	// 配置
	do.Provide(injector, services.NewConfService)

	// 日志
	do.Provide(injector, services.NewLogger)

	// grpc
	do.Provide(injector, services.NewGrpcServer)
	do.Provide(injector, services.NewApisHelloClientService)

	// gin
	do.Provide(injector, services.NewGinService)

	// 启动
	server := do.MustInvoke[*services.GinService](injector)
	server.Run()
}
