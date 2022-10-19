package trial

import (
	"os"

	"github.com/rs/zerolog"
)

func NewZLogger() (*zerolog.Logger, error) {
	// 命令行输出
	cw := zerolog.ConsoleWriter{Out: os.Stdout}

	// 文件输出 TODO 这个要找个能文件切割的
	fw, err := os.OpenFile("debug.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// 多输出端
	m := zerolog.MultiLevelWriter(cw, fw)
	logger := zerolog.New(m).With().Timestamp().Logger()
	return &logger, err
}
