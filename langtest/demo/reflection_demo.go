package demo

import (
	"fmt"
)

// type -> interface
//      -> struct  -> *pointer
//

type IAA interface {
	Faa()
}

type IBB interface {
	IAA
	Fbb()
}

type SAA struct {
	A int
}

func (i *SAA) Faa() {
	fmt.Println("i Saa")
}

type SBB struct {
	SAA
}

func (i *SBB) Fbb() {
	fmt.Println("i Sbb")
}

// 输出：
//start
//SAA {1}
//IBB &{{2}}

func ReflectionDemo() {
	fmt.Println("start")

	var aa interface{}
	aa = SAA{A: 1}

	var aa2 IBB
	aa2 = &SBB{SAA: SAA{A: 2}}

	switch v := aa.(type) {
	case SAA:
		fmt.Println("SAA", v)
	default:
		fmt.Println("else1")
	}

	switch v := aa2.(type) {
	case IBB:
		fmt.Println("IBB", v)
	case IAA:
		fmt.Println("IAA", v)
	default:
		fmt.Println("else2")
	}
}
