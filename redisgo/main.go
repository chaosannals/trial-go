package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var (
	redisPool   *redis.Pool
	redisServer = flag.String("redisServer", ":6379", "")
)

func init() {
	redisPool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", *redisServer)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}

func RedisSet(key string, v string) error {
	conn := redisPool.Get()
	_, err := conn.Do("set", key, v)
	return err
}

func RedisGet(key string) (string, error) {
	conn := redisPool.Get()
	v, err := conn.Do("get", key)
	if err != nil {
		return "", err
	}
	return redis.String(v, nil)
}


func main() {
	if err := RedisSet("aaa", "vvvvv"); err != nil {
		fmt.Println(err)
	}
	v, err := RedisGet("aaa")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
}