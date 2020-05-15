package freja

import (
	"context"
)

//go:generate moq -out middleware_mock_test.go . Middleware
// Middleware is the manager for any task that is meant to run in the background besides the http server
// It'll initiate it, run it in a different goroutine and if the service is shutdown it'll stop it gracefully.
// It can be also used as a mean to ensure the correct clean up of resources before shutdown,
// such as closing DB connections, etc.
// --- Init() ---
// This method is for any initial instructions the task may need, please bear in mind this is mean to run fast
// and it does run in the main thread, if an error is thrown, it'll prevent the service from starting
// --- Run(ctx) ---
// This method runs in a different go routine, and that process can run for as long as its needed.
// If the service gets a shutdown signal it'll cancel the context first before running the Stop(ctx) method.
// Its worth mention that if the execution of this method fails it won't re-start the service, it'll log it, if you wish
// to restart it, consider using the HealthChecker interface.
// --- Stop(ctx) ----
// This method is executed as the last step of the graceful shutdown and is meant to clean up resources before shut down
// One last thing to note, if any Middleware implements the HealthChecker interface, it'll be added automatically
// to the health check pool when added to the App
type Middleware interface {
	Init() error
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}
