package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func Max[T ~int | uint](a T, b T) T {
	if a > b {
		return a
	}
	return b
}

func TypeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func ValueAs[T any](v reflect.Value) (r T) {
	reflect.ValueOf(&r).Elem().Set(v)
	return
}

func JsonAs[T interface{}](text string) (T, error){
	var r T
	// json.NewDecoder().Decode(&r)
	err := json.Unmarshal([]byte(text), &r)
	// fmt.Println(text, r)
	return r, err
}

func De[T interface{}](text string) T {
	var r T
	t := reflect.TypeOf(r)
	// 必须通过这样才能对原来的对象字段进行操作 ValueOf(r) 不行。
	tv := reflect.ValueOf(&r).Elem()
	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := tv.Field(i)
		fvt := fv.Type()
		fmt.Println(ft.Name, "=>", fvt)
		if fv.CanSet() {
			switch fvt.Kind() {
			case reflect.Int:
				fv.SetInt(1234)
				break
			case reflect.Float32:
				fv.SetFloat(123.4)
				break
			}
		}
	}
	fmt.Print("r", r)
	return r
}