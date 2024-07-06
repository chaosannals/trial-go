package demo

import (
	"context"
	"fmt"
	"time"
)

func goLv2(i int) {
	for {
		now := time.Now()
		fmt.Printf("goLv2:[%d] %s\n", i, now.Format("2006-01-02"))
		time.Sleep(1 * time.Second)
	}
}

func DemoCancel() {
	//ctx 没有传染性，只能关闭被 select Done 的单层，子 goLv2 不受影响。
	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		i := 0
		for {
			select {
			case <-ctx.Done():
				return
			default:
				now := time.Now()
				fmt.Printf("tick: %s\n", now.Format("2006-01-02"))
			}

			time.Sleep(1 * time.Second)
			go goLv2(i)
			i++
		}
	}(ctx)

	for i := 0; i < 20; i++ {
		now := time.Now()
		time.Sleep(1 * time.Second)
		fmt.Printf("main: %s\n", now.Format("2006-01-02"))
	}
	cancel()
	for i := 0; i < 20; i++ {
		now := time.Now()
		time.Sleep(1 * time.Second)
		fmt.Printf("main2: %s\n", now.Format("2006-01-02"))
	}
}
