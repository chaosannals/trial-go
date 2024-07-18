package main

import (
	"fmt"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
)

type dvs struct {
	i int
	d float64
	s string
}

func main() {
	input := []any{"Samuel", "John", "Samuel", 123, 123.44, dvs{i: 12, d: 12.3, s: "oooo"}}
	names := lo.Map(
		lo.Uniq(input),
		func(v any, i int) string {
			return fmt.Sprintf("map[%d]: %v", i, v)
		},
	)
	fmt.Printf("unique: %v\n", names)

	names = lo.Shuffle(names)

	result := lo.Reduce(names, func(agg any, v string, i int) any {
		return fmt.Sprintf("%v + %v[%d]", agg, v, i)
	}, "0")

	fmt.Printf("reduce: %v\n", result)

	lmr := lop.Map(input, func(v any, i int) string {
		return fmt.Sprintf("map[%d]: %v", i, v)
	})
	fmt.Printf("p map: %v\n", lmr)
}
