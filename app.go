package freya

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/diego1q2w/freya/component"
	"github.com/diego1q2w/freya/healthcheck"
	"github.com/diego1q2w/freya/middleware"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:generate moq -out middleware_mock_test.go -pkg freya ./middleware Middleware

//go:generate moq -out health_calculator_mock_test.go . healthCalculator
type healthCalculator interface {
	Add(healthcheck.HealthChecker)
	Calculate() (bool, []healthcheck.Status)
}

//go:generate moq -out server_mock_test.go . Server
type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type OptionApp func(*app)

type app struct {
	healthCalculator healthCalculator
	logger           Logger
	meddlers         []middleware.Middleware
	server           Server
	cancel           context.CancelFunc
	osSignal         chan os.Signal //  listen when the service is asked to shutdown and initiates graceful shutdown
	shutdown         chan bool      // Use to exit the method start
}

func App(options ...OptionApp) *app {
	app := &app{
		healthCalculator: healthcheck.NewHealthCalculator(),
		logger:           component.NewLogger(),
		meddlers:         make([]middleware.Middleware, 0),
		osSignal:         make(chan os.Signal, 1),
		shutdown:         make(chan bool, 1),
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

func optionLogger(logger Logger) OptionApp {
	return func(a *app) {
		a.logger = logger
	}
}

// AddMiddleware adds another middleware, and if it does implement the interface HealthChecker
// it'll add as a HealCheck as well
func (a *app) AddMiddleware(m middleware.Middleware) {
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
	a.server = component.NewServer(h)
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

	signal.Notify(a.osSignal,
		os.Interrupt,    // interrupt is syscall.SIGINT, Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)

	go a.stop(context.Background())
	return nil
}

func (a *app) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	if err := a.init(); err != nil {
		return err
	}

	for _, mid := range a.meddlers {
		go func(mid middleware.Middleware) {
			if err := mid.Run(ctx); err != nil {
				a.logger.Errorf("unable to run middleware: %s", err)
			}
		}(mid)
	}

	time.Sleep(time.Millisecond * 200)

	if a.server != nil {
		if err := a.server.ListenAndServe(); err != nil {
			return fmt.Errorf("unable to run the server: %w", err)
		}
	}

	<-a.shutdown
	a.logger.Info("shutdown finalized")
	return nil
}

func (a *app) stop(ctx context.Context) {
	<-a.osSignal
	a.logger.Info("shutdown initiated")
	a.cancel()

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	if a.server != nil {
		if err := a.server.Shutdown(ctx); err != nil {
			a.logger.Errorf("error stopping server: %s", err)
		}
	}

	for _, mid := range a.meddlers {
		if err := mid.Stop(ctx); err != nil {
			a.logger.Errorf("error stopping middleware: %s", err)
		}
	}

	a.shutdown <- true
}
