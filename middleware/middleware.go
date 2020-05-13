package middleware

import "context"

type MiddleWare interface {
	Init() error
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
	Name() string
}
