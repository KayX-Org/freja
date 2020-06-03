package middleware

import (
	"context"
	"fmt"
	"github.com/kayx-org/freja/component"
	"github.com/kayx-org/freja/healthcheck"
	"net"
)

type grpcMiddleware struct {
	service  *component.GRPCServer
	listener *net.Listener
	status   healthcheck.ServiceStatus
}

func NewGRPCMiddleware(service *component.GRPCServer) *grpcMiddleware {
	return &grpcMiddleware{
		service: service,
		status:  healthcheck.UP,
	}
}

func (m *grpcMiddleware) Init() error {
	lis, err := net.Listen("tcp", m.service.ListenAddress())
	if err != nil {
		return fmt.Errorf("unable to bind listener: %w", err)
	}
	m.listener = &lis

	return nil
}

func (m *grpcMiddleware) Run(context.Context) error {
	defer func() {
		m.status = healthcheck.DOWN
	}()

	if err := m.service.Server().Serve(*m.listener); err != nil {
		return fmt.Errorf("unabel to serve: %w", err)
	}

	return nil
}

func (m *grpcMiddleware) Stop(context.Context) error {
	m.service.Server().GracefulStop()
	return nil
}

func (m *grpcMiddleware) Name() string {
	return "grpc"
}

func (m *grpcMiddleware) Status() healthcheck.ServiceStatus {
	return m.status
}
