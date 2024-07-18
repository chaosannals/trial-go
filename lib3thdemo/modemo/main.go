package main

import (
	"fmt"

	"github.com/samber/mo"
)

func main() {
	option1 := mo.Some(42)
	r1 := option1.Map(func(value int) (int, bool) {
		return value + 100, false
	})
	fmt.Printf("r1: %v\n", r1)
	r2 := r1.Match(func(value int) (int, bool) {
		return 1, true
	}, func() (int, bool) {
		return 0, true
	})
	fmt.Printf("r2: %v\n", r2)
}
