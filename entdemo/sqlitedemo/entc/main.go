package main

import (
	"log"
	"net/http"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func main() {
	opts := []entc.Option{
        entc.Dependency(
            entc.DependencyType(&http.Client{}),
        ),
        entc.Dependency(
            entc.DependencyName("Writer"),
            entc.DependencyTypeInfo(&field.TypeInfo{
                Ident:   "io.Writer",
                PkgPath: "io",
            }),
        ),
    }
	if err := entc.Generate("../ent/schema", &gen.Config{}, opts...); err != nil {
		log.Fatal("running ent codegen:", err)
	}
}
