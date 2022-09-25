package main

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"trial/wstpd/basis"
	"trial/wstpd/frontend"
	"trial/wstpd/software"
)

// 初始化
func init() {
	logf, err := rotatelogs.New(
		"runtime/logs/%Y%m%d.log",
		rotatelogs.WithMaxAge(15*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(logf)

	// 加载环境变量
	err = godotenv.Load("runtime/.env")
	if err != nil {
		log.Println(".env 未加载")
	} else {
		log.Println(".env 已加载")
	}

	// 配置日志。
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.SetFormatter(&log.TextFormatter{})
}

// 主入口
func main() {
	// 初始化
	server, err := basis.NewBizServer()
	if err != nil {
		log.Error(err)
		return
	}

	// 前端
	f, err := frontend.NewFrontendDispatcher(server)
	if err != nil {
		log.Error(err)
		return
	}
	server.Attach("/biz", f.DispatchRequest)

	// 服务
	s, err := software.NewSoftwareDispatcher(server)
	if err != nil {
		log.Error(err)
		return
	}
	server.Attach("/svc", s.DispatchRequest)

	// 启动
	server.Run()
}
