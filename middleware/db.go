package middleware

import (
	"context"
	"database/sql"
)

type dbMiddleware struct {
	db *sql.DB
}

func NewDB(db *sql.DB) *dbMiddleware {
	return &dbMiddleware{db: db}
}

func (m *dbMiddleware) Init() error {
	return nil
}

func (m *dbMiddleware) Run(ctx context.Context) error {
	return nil
}

func (m *dbMiddleware) Stop(ctx context.Context) error {
	return m.db.Close()
}
