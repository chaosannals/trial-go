package app

import (
	"fmt"

	"gopkg.in/natefinch/lumberjack.v2"
)

type DemoLogger struct {
	logger *lumberjack.Logger
}

func NewDemoLogger() *DemoLogger {
	logger := lumberjack.Logger{
		Filename:   "./log/myapp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}
	return &DemoLogger{
		logger: &logger,
	}
}

func (dl *DemoLogger) Printf(format string, v ...interface{}) {
	text := fmt.Sprintf(format, v...)
	dl.logger.Write([]byte(text))
}
