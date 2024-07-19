package demo

import (
	"fmt"
)

type YetA string
type YetB string

type YetS[T any] struct {
	Value T
}

func see(vs ...any) {
	for i, v := range vs {
		fmt.Printf("%d %v", i, v)
		switch v.(type) {
		case YetS[YetA]:
			fmt.Println("yet A.")
		case YetS[YetB]:
			fmt.Println("yet B.")
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
	c := YetS[string]{}
	d := YetS[string]{}
	see(a, b, c, d)
}
