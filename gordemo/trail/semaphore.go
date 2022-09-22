package trail

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

const (
	semaphoreCount     = 100
	semaphoreTaskCount = 1000000
	semaphoreWeight    = 1
)

func TrailSemaphore() {
	sem := semaphore.NewWeighted(semaphoreCount)
	for i := 0; i < semaphoreTaskCount; i += 1 {
		if err := sem.Acquire(context.Background(), semaphoreWeight); err != nil {
			fmt.Printf("sem acquire err: %v\n", err)
			break
		}

		go func(i int) {
			defer sem.Release(semaphoreWeight)
			// 死循环可以拉满
			//for { }

			// 10C20T 跑不满 IO 打印 IO 导致
			fmt.Printf("go %d \n", i)
		}(i)
	}

	if err := sem.Acquire(context.Background(), semaphoreCount); err != nil {
		fmt.Printf("final sem acquire err: %v\n", err)
	}

	fmt.Println("ending--------")
}
