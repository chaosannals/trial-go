package demo

import "fmt"

func HasRecover(name string, action func()) {
	fmt.Printf("[%s] start\n", name)
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[%s] recover: %v\n", name, r)
		}
	}()
	action()
	fmt.Printf("[%s] end\n", name)
}

func PanicRecoverDemo() {
	fmt.Println("start")
	HasRecover("demo 1", func() {
		panic("I's panic.")
	})
	fmt.Println("step 1")
	HasRecover("demo 2", func() {
		panic(map[string]any{
			"a": 123,
			"b": "asdf",
		})
	})
	fmt.Println("step 2")
	HasRecover("demo 3", func() {
		for i := 10; i >= 0; i-- {
			fmt.Println(100 / i) // 除0 可以被 recover
		}
	})
	fmt.Println("end")
}
