package demo

import "fmt"

// type = 号的差距应该是泛型加入后的事
// 没有 = 号被认作不同类型，仅放射和泛型时会被识别，一般操作会隐式自动转换类型。
type S1 = string
type N1 string

type SS1[T any] struct {
	v T
}

const (
	S1_V1 S1 = "s1"
	N1_V1 N1 = "s1"
)

func TypeDefineDemo() {
	b1 := S1_V1 == "s1"
	fmt.Printf("b1: %v\n", b1)

	b2 := N1_V1 == "s1" // 不同类型，居然可以比较。
	fmt.Printf("b2: %v\n", b2)

	ss1v1 := SS1[string]{}
	ss1v2 := SS1[S1]{}
	// ss1v3 := SS1[N1]{}

	b3 := ss1v1 == ss1v2
	fmt.Printf("b3: %v\n", b3)
	// b4 := ss1v1 == ss1v3 // 不同类型
}
