package main

import (
	"fmt"
)

func main() {
	// 没有长度的 chan 读写必须不同协程
	// var ch1 chan string = make(chan string)
	// 直接会被检测出 fatal error: all goroutines are asleep - deadlock!

	// 管道带了数量后，即使是 1 也和没有数量的不同。
	var ch1 chan string = make(chan string, 1)
	// time.Sleep(time.Second * 2)
	ch1 <- "result 1"
	fmt.Println("ch1 result:", <-ch1)
}
