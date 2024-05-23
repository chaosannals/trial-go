package main

import (
	"fmt"
)

// defer 是 [函数作用域] 级别的，出 [块作用域] 不会调用。
// 要等整个函数调用结束才会被调用。
func main() {
	fmt.Println("start")
	defer func() { fmt.Println("Lv1") }()

	for i := 0; i < 10; i++ {
		fmt.Println("start: Loop Lv1")
		d := i
		defer func() { fmt.Printf("Loop Lv1: %d\n", d) }()

		for j := 0; j < 10; j++ {
			fmt.Println("start: Loop Lv2")
			v := j
			defer func() { fmt.Printf("Loop Lv2: %d\n", v) }()
			fmt.Println("end: Loop Lv2")
		}
		fmt.Println("end: Loop Lv1")
	}

	fmt.Println("end")
}
