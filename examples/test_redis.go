package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})
	
	res := rdb.HSet(ctx, "2020_09_16_15_14", "v9111", 10086)
	fmt.Println(res.Result())
}
