package main

import (
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"fmt"
	"syscall/js"
)

func Export(name string, module map[string]interface{}) {
	js.Global().Set(name, module)
}

func MakeJsFunc() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		for i := 0; i < len(args); i += 1 {
			fmt.Println(args[i].String())
		}
		return this.String()
	})
}

// 导出函数，必须常驻。
func ExportHold() {
	fmt.Println("开始加载")
	Export("demo", map[string]interface{}{
		// 函数的格式固定是 func(js.Value, []js.Value) interface{} ;
		"hello": MakeJsFunc(),
	})
	fmt.Println("加载完成")
	select {} // 必须常驻，保证导出函数有效。
}

// 遍历 JS 数组。
func IterArray(a js.Value) {
	for i := 0; i < a.Length(); i += 1 {
		fmt.Println(a.Index(i).String())
	}
}

// 具体查看 go 源码
// src/syscall/js 下定义
func main() {
	ExportHold()
}
