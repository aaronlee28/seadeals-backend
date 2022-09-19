package redis

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

var rds *redis.Client

func Setup() {
	rds = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	fmt.Println("redis client initialized")
}

func Use() *redis.Client {
	return rds
}

// REDIS EXAMPLE

// redis.Setup()
// ctx := context.Background()
// fmt.Println(redis.Use().Get(ctx, "test"))
