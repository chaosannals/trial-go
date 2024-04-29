package main

import "fmt"

func main() {
	ca := make(chan int, 1)
	cb := make(chan int, 1)

	go func() {
		for i := 0; i < 100; i += 1 {
			fmt.Print("A")
			ca <- i
			<-cb
		}
		close(ca)
	}()

	for range ca {
		fmt.Println("B")
		cb <- 1
	}
}
