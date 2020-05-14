package component

import (
	"fmt"
	"net/http"
	"time"
)

func NewServer(handler http.Handler) *http.Server {
	addr := getEnv("SERVICE_ADDR", "0.0.0.0")
	port := getEnv("SERVICE_PORT", "5042")

	return &http.Server{
		Addr:        fmt.Sprintf("%s:%s", addr, port),
		IdleTimeout: time.Second * 2,
		Handler:     handler,
	}
}
