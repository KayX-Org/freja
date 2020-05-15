package middleware

import (
	"context"
	"github.com/diego1q2w/freja/healthcheck"
	"time"
)

//go:generate moq -out db_mock_test.go . db
type db interface {
	Close() error
	PingContext(context.Context) error
}

type OptionDbMiddleware func(*dbMiddleware)

type dbMiddleware struct {
	db          db
	name        string
	checkWindow time.Duration
	status      healthcheck.ServiceStatus
}

// NewDB returns a new DB middleware which also implements the HealthCheck interface and can be configured accordinginly
func NewDB(db db, options ...OptionDbMiddleware) *dbMiddleware {
	midDb := &dbMiddleware{
		db:          db,
		name:        "db",
		checkWindow: time.Millisecond * 200,
		status:      healthcheck.UP,
	}

	for _, op := range options {
		op(midDb)
	}

	return midDb
}

func OptionWindowCheck(t time.Duration) OptionDbMiddleware {
	return func(m *dbMiddleware) {
		m.checkWindow = t
	}
}

func OptionHealthCheckName(t string) OptionDbMiddleware {
	return func(m *dbMiddleware) {
		m.name = t
	}
}

func (m *dbMiddleware) Init() error {
	return nil
}

func (m *dbMiddleware) Run(ctx context.Context) error {
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

func (m *dbMiddleware) Stop(context.Context) error {
	return m.db.Close()
}

func (m *dbMiddleware) Name() string {
	return m.name
}

func (m *dbMiddleware) runStatusCheck(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if err := m.db.PingContext(ctx); err != nil {
		m.status = healthcheck.DOWN
	} else {
		m.status = healthcheck.UP
	}
}

func (m *dbMiddleware) Status() healthcheck.ServiceStatus {
	return m.status
}
