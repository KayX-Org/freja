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
	httpSrv *http.Server
	logger  io.Writer
}

func NewServer(handler http.Handler, options ...OptionServer) *Server {
	srv := &Server{}
	for _, o := range options {
		o(srv)
	}

	addr := getEnv("SERVICE_ADDR", "0.0.0.0")
	port := getEnv("SERVICE_PORT", "5042")
	srv.httpSrv = &http.Server{
		Addr:        fmt.Sprintf("%s:%s", addr, port),
		IdleTimeout: time.Second * 2,
		Handler:     handler,
	}
	if srv.logger != nil {
		srv.httpSrv.ErrorLog = log.New(srv.logger, "", 0)
	}

	srv.httpSrv.SetKeepAlivesEnabled(true)

	return srv
}

func OptionErrorLogWriter(log io.Writer) OptionServer {
	return func(server *Server) {
		server.logger = log
	}
}

func (s *Server) ListenAndServe() error {
	return s.httpSrv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
