package freja

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kayx-org/freja/component"
	"github.com/kayx-org/freja/healthcheck"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:generate moq -out health_calculator_mock_test.go . healthCalculator
type healthCalculator interface {
	Add(healthcheck.HealthChecker)
	Calculate() (bool, []Status)
}

//go:generate moq -out server_mock_test.go . Server
type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type OptionApp func(*app)

type app struct {
	healthCalculator        healthCalculator
	gracefulShutdownTimeout time.Duration
	logger                  Logger
	meddlers                []Middleware
	server                  Server
	cancel                  context.CancelFunc
	osSignal                chan os.Signal //  listen when the service is asked to shutdown
	gracefulStop            chan bool      // Use to initiate the graceful shutdown
}

func App(options ...OptionApp) *app {
	app := &app{
		healthCalculator:        NewHealthCalculator(),
		logger:                  component.NewLogger(),
		gracefulShutdownTimeout: time.Second * 10,
		meddlers:                make([]Middleware, 0),
		osSignal:                make(chan os.Signal, 1),
		gracefulStop:            make(chan bool, 1),
	}

	for _, o := range options {
		o(app)
	}

	return app
}

func OptionHealthCalculator(healthCalculator healthCalculator) OptionApp {
	return func(a *app) {
		a.healthCalculator = healthCalculator
	}
}

func OptionGracefulShutdownTimeout(gracefulShutdownTimeout time.Duration) OptionApp {
	return func(a *app) {
		a.gracefulShutdownTimeout = gracefulShutdownTimeout
	}
}

func optionLogger(logger Logger) OptionApp {
	return func(a *app) {
		a.logger = logger
	}
}

// AddMiddleware adds another middleware, and if it does implement the interface HealthChecker
// it'll add as a HealCheck as well
func (a *app) AddMiddleware(m Middleware) {
	a.meddlers = append(a.meddlers, m)
	if h, ok := m.(healthcheck.HealthChecker); ok {
		a.AddHealthCheck(h)
	}
}

func (a *app) AddHealthCheck(h healthcheck.HealthChecker) {
	if a.healthCalculator != nil {
		a.healthCalculator.Add(h)
	}
}

// HealthCheck returns a boolean whether the service is healthy or not, and also accepts an the io.Writer
// into which writes the summary of all the health checks in JSON format
func (a *app) HealthCheck(writer io.Writer) (bool, error) {
	if a.healthCalculator != nil {
		status, summary := a.healthCalculator.Calculate()
		if err := json.NewEncoder(writer).Encode(summary); err != nil {
			return false, fmt.Errorf("unable to encode the summary: %w", err)
		}
		return status, nil
	}
	return true, nil
}

func (a *app) WithServer(s Server) {
	a.server = s
}

func (a *app) Server(h http.Handler) {
	a.server = component.NewServer(h,
		component.OptionErrorLogWriter(
			NewLogWrite(a.logger, "error")))
}

func (a *app) Logger() Logger {
	return a.logger
}

func (a *app) init() error {
	for _, mid := range a.meddlers {
		if err := mid.Init(); err != nil {
			return fmt.Errorf("unable to run Init(): %w", err)
		}
	}

	signal.Notify(a.osSignal, syscall.SIGTERM)
	signal.Notify(a.osSignal, syscall.SIGINT)

	go func() {
		sig := <-a.osSignal
		a.logger.Infof("signal caught: %s", sig)
		a.gracefulStop <- true
	}()

	return nil
}

func (a *app) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	if err := a.init(); err != nil {
		return err
	}

	for _, mid := range a.meddlers {
		go func(mid Middleware) {
			if err := mid.Run(ctx); err != nil {
				a.logger.Errorf("unable to run middleware: %s", err)
			}
		}(mid)
	}

	time.Sleep(time.Millisecond * 200)

	if a.server != nil {
		go func() {
			a.logger.Info("starting server")
			if err := a.server.ListenAndServe(); err != nil {
				a.logger.Fatalf("unable to run the server: %s", err)
			}
		}()

	}

	a.stop(ctx)

	a.logger.Info("shutdown finalized")
	return nil
}

func (a *app) stop(ctx context.Context) {
	<-a.gracefulStop
	a.logger.Info("graceful shutdown initiated")
	a.cancel()

	ctx, cancel := context.WithTimeout(ctx, a.gracefulShutdownTimeout)
	defer cancel()

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Errorf("error gracefully stopping server: %s", err)
		}
	}

	for _, mid := range a.meddlers {
		if err := mid.Stop(ctx); err != nil {
			a.logger.Errorf("error gracefully stopping middleware: %s", err)
		}
	}
}
