package freya

import (
	"github.com/diego1q2w/freya/healthcheck"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	testCases := map[string]struct {
		healthChecks    []*mockHC
		expectedStatus  bool
		expectedSummary []Status
	}{
		"if all the health check is empty true expected": {
			healthChecks:    []*mockHC{},
			expectedStatus:  true,
			expectedSummary: []Status{},
		},
		"if all the health checks are correct or temporally unavailable then true expected": {
			healthChecks: []*mockHC{
				{name: "foo", status: healthcheck.UP},
				{name: "bar", status: healthcheck.TemporallyUnavailable},
			},
			expectedStatus: true,
			expectedSummary: []Status{
				{Name: "foo", Status: healthcheck.UP.ToString()},
				{Name: "bar", Status: healthcheck.TemporallyUnavailable.ToString()},
			},
		},
		"if one  health checks is down then false expected": {
			healthChecks: []*mockHC{
				{name: "foo", status: healthcheck.UP},
				{name: "bar", status: healthcheck.DOWN},
			},
			expectedStatus: false,
			expectedSummary: []Status{
				{Name: "foo", Status: healthcheck.UP.ToString()},
				{Name: "bar", Status: healthcheck.DOWN.ToString()},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			hc := NewHealthCalculator()

			for _, h := range tc.healthChecks {
				hc.Add(h)
			}

			status, summary := hc.Calculate()
			assert.Equal(t, tc.expectedStatus, status)
			assert.Equal(t, tc.expectedSummary, summary)
		})
	}
}

type mockHC struct {
	name   string
	status healthcheck.ServiceStatus
}

func (m *mockHC) Name() string {
	return m.name
}
func (m *mockHC) Status() healthcheck.ServiceStatus {
	return m.status
}
