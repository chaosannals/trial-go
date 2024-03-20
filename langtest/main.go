package main

import (
	"fmt"
)

func main() {
	fmt.Println("start")

	// vv := []byte{12, 34, 54}
	// a, isOk := vv[6]
	// fmt.Println(a)

	mm := map[string]string{
		"aa": "1232",
		"bb": "234324",
	}

	v, isOk := mm["12"]
	fmt.Println(v)
	fmt.Println(isOk)

	cc := make(chan string, 1)

	go func() {
		cc <- "sssss"
		close(cc)
	}()

	for rr := range cc {
		fmt.Println(rr)
	}
	// fmt.Println(isOk)
}
