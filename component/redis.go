package component

import (
	"github.com/go-redis/redis/v8"
)

func NewClientRedis(db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       getEnv("REDIS_ADDR", "localhost:6379"),
		Username:   getEnv("REDIS_USERNAME", ""),
		Password:   getEnv("REDIS_PASSWORD", ""), // no password set
		DB:         db,
		MaxRetries: getEnvAsInt("REDIS_MAX_RETRIES", 3),
	})
}
