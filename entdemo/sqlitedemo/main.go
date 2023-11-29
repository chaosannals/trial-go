package main

import (
	"context"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/chaosannals/sqlitedemo/ent"
	"github.com/chaosannals/sqlitedemo/ent/hook"
)

func main() {
	client, err := ent.Open(
		"sqlite3",
		// "file:ent?mode=memory&cache=shared&_fk=1",
		"file:./demo.db?cache=shared&_fk=1",
		// ent.Writer(os.Stdout),
		// ent.HTTPClient(http.DefaultClient),
	)
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	// 在生成构建器中使用注入依赖项的示例。
	client.User.Use(func(next ent.Mutator) ent.Mutator {
		return hook.UserFunc(func(ctx context.Context, m *ent.UserMutation) (ent.Value, error) {
			// _ = m.HTTPClient
			// _ = m.Writer
			return next.Mutate(ctx, m)
		})
	})

	ctx := context.Background()

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	u1, err := client.User.
		Create().
		Save(ctx)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(u1)

	if err := client.User.DeleteOneID(u1.ID).Exec(ctx); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("delete")
	}
}
