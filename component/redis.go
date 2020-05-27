package component

import (
	"github.com/go-redis/redis/v8"
)

func NewClientRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", ""), // no password set
		DB:       getEnvAsInt("REDIS_DATABASE", 0),
	})
}
