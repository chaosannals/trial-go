package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ctx context.Context
	rdb redis.Client
)


func init() {
	ctx = context.Background()
	rdb = *redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
}

func main() {
	if err := rdb.Set(ctx, "tkey", "vvvvv", time.Second * 60).Err(); err != nil {
		panic(err)
	}
	if v, err := rdb.Get(ctx, "tkey").Result(); err != nil {
		panic(err)
	} else {
		fmt.Println("key", v)
	}
	if vn, err := rdb.Get(ctx, "nkey").Result(); err == redis.Nil {
		fmt.Println("nkey 不存在")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("nkey", vn)
	}
}