package freja

import "github.com/kayx-org/freja/healthcheck"

type Status struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type healthCalculate struct {
	healthChecks []healthcheck.HealthChecker
}

func NewHealthCalculator() *healthCalculate {
	return &healthCalculate{healthChecks: make([]healthcheck.HealthChecker, 0)}
}

func (h *healthCalculate) Add(healthCheck healthcheck.HealthChecker) {
	h.healthChecks = append(h.healthChecks, healthCheck)
}

func (h *healthCalculate) Calculate() (bool, []Status) {
	statuses := make([]Status, 0)
	finalStatus := true

	for _, hc := range h.healthChecks {
		if hc.Status().IsDown() {
			finalStatus = false
		}

		statuses = append(statuses, Status{
			Name:   hc.Name(),
			Status: hc.Status().ToString(),
		})
	}

	return finalStatus, statuses
}
