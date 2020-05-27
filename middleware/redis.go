package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/kayx-org/freja/healthcheck"
	"time"
)

type redisClient interface {
	Ping(ctx context.Context) *redis.StatusCmd
	Close() error
}

type redisMiddleware struct {
	client      redisClient
	name        string
	checkWindow time.Duration
	status      healthcheck.ServiceStatus
}

func NewRedisMiddleware(client redisClient) *redisMiddleware {
	return &redisMiddleware{
		client:      client,
		name:        "redis",
		checkWindow: time.Second,
		status:      healthcheck.UP,
	}
}

func (m *redisMiddleware) Init() error {
	_, err := m.client.Ping(context.Background()).Result()
	return err
}

func (m *redisMiddleware) Run(ctx context.Context) error {
	ticker := time.NewTicker(m.checkWindow)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			m.runStatusCheck(ctx)
		}
	}
}

func (m *redisMiddleware) Stop(context.Context) error {
	return m.client.Close()
}

func (m *redisMiddleware) Name() string {
	return m.name
}

func (m *redisMiddleware) runStatusCheck(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, m.checkWindow)
	defer cancel()

	if _, err := m.client.Ping(ctx).Result(); err != nil {
		m.status = healthcheck.DOWN
	} else {
		m.status = healthcheck.UP
	}
}

func (m *redisMiddleware) Status() healthcheck.ServiceStatus {
	return m.status
}
