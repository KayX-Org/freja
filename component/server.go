package component

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type OptionServer func(*Server)

type Server struct {
	srv    *http.Server
	logger io.Writer
}

func NewServer(handler http.Handler, options ...OptionServer) *Server {
	srv := &Server{}
	for _, o := range options {
		o(srv)
	}

	addr := getEnv("SERVICE_ADDR", "0.0.0.0")
	port := getEnv("SERVICE_PORT", "5042")
	httpServer := &http.Server{
		Addr:        fmt.Sprintf("%s:%s", addr, port),
		IdleTimeout: time.Second * 2,
		Handler:     handler,
	}
	if srv.logger != nil {
		httpServer.ErrorLog = log.New(srv.logger, "", 0)
	}

	httpServer.SetKeepAlivesEnabled(true)

	return srv
}

func OptionErrorLogWriter(log io.Writer) OptionServer {
	return func(server *Server) {
		server.logger = log
	}
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
