package repository

import "github.com/go-redis/redis/v8"

// Redis client
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis address
	})
}
