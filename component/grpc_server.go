package component

import (
	"fmt"
	"google.golang.org/grpc"
	"time"
)

type GRPCServer struct {
	port   string
	addr   string
	server *grpc.Server
}

func NewGRPCServer() *GRPCServer {
	addr := getEnv("GRPC_SERVICE_ADDR", "0.0.0.0")
	port := getEnv("GRPC_SERVICE_PORT", "50051")

	s := grpc.NewServer(grpc.ConnectionTimeout(time.Second * 10))

	return &GRPCServer{
		addr:   addr,
		port:   port,
		server: s,
	}
}

func (s *GRPCServer) Server() *grpc.Server {
	return s.server
}

func (s *GRPCServer) ListenAddress() string {
	return fmt.Sprintf("%s:%s", s.addr, s.port)
}
