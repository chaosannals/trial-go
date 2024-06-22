package demo

import "fmt"

// panic defer recover 机制比其他语言 的 try catch 机制，作用域上不是随意的。
// 因为 defer 是函数作用域，所以 recover 的作用域也只能是函数作用域
// 其他语言的 try catch 可以直接包裹代码形成作用域
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
