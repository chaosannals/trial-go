package demo

import (
	"fmt"
)

// GOLANG 在此种又特别严格。此处 YetA 和 YetB 不是 string 别名这么简单。而是实打实的另外一个类型。
type YetA string
type YetB string

// 加上等于号后，NameA 和 NameB 只是别名，本质上等价于 string
type NameA = string
type NameB = string

type YetS[T any] struct {
	Value T
}

// GOLANG 在此种又特别严格。此处 YetSA 和 YetSB 不是 YetS[YetA] YetS[YetB] 别名这么简单。而是实打实的另外一个类型。
type YetSA YetS[YetA]
type YetSB YetS[YetB]

// 此处正是别名
type NameSA = YetS[YetA]
type NameSB = YetS[YetB]

func see(vs ...any) {
	for i, v := range vs {
		fmt.Printf("%d %v", i, v)
		switch v.(type) {
		case YetS[YetA]:
			fmt.Println("yet A.")
		case YetS[YetB]:
			fmt.Println("yet B.")
		case YetSA:
			fmt.Println("yet SA.")
		case YetSB:
			fmt.Println("yet SB.")
		case YetS[string]:
			fmt.Println("yet string.")
		default:
			fmt.Println("yet other.")
		}
	}
}

func ReflectionTypenameDemo() {
	a := YetS[YetA]{}
	b := YetS[YetB]{}
	sa := YetSA{}
	sb := YetSB{}
	na := YetS[NameA]{}
	nb := YetS[NameB]{}
	nsa := NameSA{}
	nsb := NameSB{}
	c := YetS[string]{}
	d := YetS[string]{}
	see(a, b, sa, sb, na, nb, nsa, nsb, c, d)
}
