package healthcheck

import (
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
				{name: "foo", status: UP},
				{name: "bar", status: TemporallyUnavailable},
			},
			expectedStatus: true,
			expectedSummary: []Status{
				{Name: "foo", Status: "up"},
				{Name: "bar", Status: TemporallyUnavailable},
			},
		},
		"if one  health checks is down then false expected": {
			healthChecks: []*mockHC{
				{name: "foo", status: UP},
				{name: "bar", status: DOWN},
			},
			expectedStatus: false,
			expectedSummary: []Status{
				{Name: "foo", Status: "up"},
				{Name: "bar", Status: DOWN},
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
	status ServiceStatus
}

func (m *mockHC) Name() string {
	return m.name
}
func (m *mockHC) Status() ServiceStatus {
	return m.status
}
