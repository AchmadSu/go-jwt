package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RedisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	DB:   0,
})
