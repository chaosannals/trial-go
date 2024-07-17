package demo

type GI1 interface {
	F1(v string) string
}

type GIDI[T any] interface {
	F2(v T) T
}

type GIDS[T any] struct {
	Gidi GIDI[T]
}

func GenericInterfaceDemo() {
	
}