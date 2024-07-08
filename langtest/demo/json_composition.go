package demo

import (
	"encoding/json"
	"fmt"
	"log"
)

// GOLANG 指针是反多态的。
// 多态只通过 接口 使用。

type JsonCompositionA struct {
	Afloat float32 `json:"afloat"`
}

type JsonCompositionB struct {
	JsonCompositionA
	Bint int32 `json:"bint"`
}

func JsonCompositionPrint[T any](v T) {
	b, err := json.Marshal(&v)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(b))
}

func JsonCompositionDemo() {
	ap := &JsonCompositionA{Afloat: 1.23}
	bp := &JsonCompositionB{
		JsonCompositionA: JsonCompositionA{Afloat: 4.56},
	}

	var i any
	i = bp
	ap = &bp.JsonCompositionA // 指针的类型会存活到运行时，此时数据还是被 any 的函数传递到反射时，还是认定是 A 类型。这种行为反多态。
	JsonCompositionPrint(&ap)
	JsonCompositionPrint(&bp)
	JsonCompositionPrint(i)
}
