package services

import (
	"log/slog"
	"os"

	"github.com/samber/do"
)

func NewLogger(i *do.Injector) (*slog.Logger, error) {
	logger := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true, // 打印全路径
			Level:     slog.LevelDebug,
		}),
	)

	return logger, nil
}
