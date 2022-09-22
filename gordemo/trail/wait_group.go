package trail

import (
	"fmt"
	"sync"
)

var (
	waitGroupChanSize       = 50
	waitGroupGoroutineCount = 100
)

//TrailWaitGroup
func TrailWaitGroup() {
	jobsChan := make(chan int, waitGroupChanSize)

	// workers
	var wg sync.WaitGroup
	for i := 0; i < waitGroupGoroutineCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range jobsChan {
				// 10C20T 跑不满 IO 打印 IO 导致
				fmt.Println(item)
			}
		}()
	}

	// senders
	for i := 0; i < 1000000; i++ {
		jobsChan <- i
	}

	// 关闭channel，上游的goroutine在读完channel的内容，就会通过wg的done退出
	close(jobsChan)
	wg.Wait()
}
