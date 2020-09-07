package middleware

import (
	"context"
	"github.com/kayx-org/freja/healthcheck"
	"github.com/kayx-org/frida"
)

type FridaMiddleware struct {
	frida *frida.Frida
	name  string
}

func NewFridaMiddleware(frida *frida.Frida) *FridaMiddleware {
	return &FridaMiddleware{frida: frida, name: "frida"}
}

func (m *FridaMiddleware) Init() error {
	return m.frida.Init()
}

func (m *FridaMiddleware) Run(ctx context.Context) error {
	return m.frida.Run(ctx)
}

func (m *FridaMiddleware) Stop(ctx context.Context) error {
	return m.frida.Stop(ctx)
}

func (m *FridaMiddleware) Name() string {
	return m.name
}

func (m *FridaMiddleware) Status() healthcheck.ServiceStatus {
	if m.frida.IsRunning() {
		return healthcheck.UP
	}

	return healthcheck.DOWN
}
