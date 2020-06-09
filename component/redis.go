package component

import (
	"github.com/go-redis/redis/v8"
	"github.com/kayx-org/freja/env"
)

func NewClientRedis(db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       env.GetEnv("REDIS_ADDR", "localhost:6379"),
		Username:   env.GetEnv("REDIS_USERNAME", ""),
		Password:   env.GetEnv("REDIS_PASSWORD", ""), // no password set
		DB:         db,
		MaxRetries: env.GetEnvAsInt("REDIS_MAX_RETRIES", 3),
	})
}
