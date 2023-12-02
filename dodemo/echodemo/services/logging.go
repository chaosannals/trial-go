package services

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/do"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewZLoggerService(i *do.Injector) (*zerolog.Logger, error) {
	// 命令行输出
	cw := zerolog.ConsoleWriter{Out: os.Stdout}

	// 简单文件输出
	fw, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// 带文件切割和压缩的输出
	lw := &lumberjack.Logger{
		Filename:   "./log/debug.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	// 多输出端
	m := zerolog.MultiLevelWriter(cw, fw, lw)
	logger := zerolog.New(m).With().Timestamp().Logger()
	return &logger, err
}
