package main

import (
	"fmt"
	"reflect"

	"github.com/chaosannals/trial-go-refl/library"
)

func main() {
	p := library.Person{}

	t := reflect.TypeOf(p)
	fmt.Println(t.Name())

	// 只有值可以列举字段。
	for i := 0; i < t.NumField(); i++ {
		fmt.Println(t.Field(i).Name)
	}

	// 使用 值，只有用 值方式传递 this 的方法。
	for i := 0; i < t.NumMethod(); i++ {
		fmt.Println(t.Method(i).Name)
	}
	m, b := t.MethodByName("SayByValue")
	if b {
		fmt.Println(m)
		a := make([]reflect.Value, 2)
		a[0] = reflect.ValueOf(p)
		a[1] = reflect.ValueOf("bbbbbb")
		m.Func.Call(a)
	}

	// ValueOf 得到 Value 类型后 Value 类型 的返回也大多是 Value 类型
	v := reflect.ValueOf(p)
	fmt.Println(v.NumMethod())
	for i := 0; i < v.NumMethod(); i++ {
		fmt.Println(v.Method(i))
	}
	vm := v.MethodByName("SayByValue")
	a := make([]reflect.Value, 1)
	a[0] = reflect.ValueOf("cccccc")
	vm.Call(a)
	fmt.Println(vm)

	fmt.Println("--------------------------")
	pp := &library.Person{}
	pt := reflect.TypeOf(pp)

	// 指针无法调用 NumField()，必须解引用。
	pvt := reflect.TypeOf(*pp)
	for i := 0; i < pvt.NumField(); i++ {
		fmt.Println(pvt.Field(i).Name)
	}

	// 使用指针，可以获得全方法信息。
	for i := 0; i < pt.NumMethod(); i++ {
		fmt.Println(pt.Method(i).Name)
	}
}
