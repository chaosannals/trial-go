package main

import (
	"fmt"

	"github.com/chaosannals/dodemo-echodemo/controllers"
	"github.com/chaosannals/dodemo-echodemo/services"
	"github.com/samber/do"
)

func main() {
	// 容器
	injector := do.New()
	defer injector.Shutdown()

	// 提供 配置
	do.Provide(injector, services.NewConfService)
	// 提供 日志服务
	do.Provide(injector, services.NewZLoggerService)
	// 提供 echo 服务
	do.Provide(injector, services.NewEchoHttpService)

	// 提供 控制器
	do.Provide(injector, controllers.NewIndexController)

	// 检查
	do.HealthCheck[services.EchoHttpService](injector)

	// 启动 http 服务
	echo, err := do.Invoke[*services.EchoHttpService](injector)
	if err != nil {
		fmt.Println(err)
	}
	echo.Serve()
}
