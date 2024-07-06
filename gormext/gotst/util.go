package main

import (
	"context"
	"time"
)

func goTimeout[F func()](timeout time.Duration, action F) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(timeout)
		select {
		case <-ctx.Done():
			return
		default:
			action()
		}
	}()

	return cancel
}
