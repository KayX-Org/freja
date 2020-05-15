package healthcheck

type ServiceStatus string

func (s ServiceStatus) ToString() string {
	return string(s)
}

func (s ServiceStatus) IsDown() bool {
	return s == DOWN
}

const (
	UP                    ServiceStatus = "up"
	DOWN                  ServiceStatus = "down"
	TemporallyUnavailable ServiceStatus = "unavailable"
)

type HealthChecker interface {
	Name() string
	Status() ServiceStatus
}

type Status struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type healthCalculator struct {
	healthChecks []HealthChecker
}

func NewHealthCalculator() *healthCalculator {
	return &healthCalculator{healthChecks: make([]HealthChecker, 0)}
}

func (h *healthCalculator) Add(healthCheck HealthChecker) {
	h.healthChecks = append(h.healthChecks, healthCheck)
}

func (h *healthCalculator) Calculate() (bool, []Status) {
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
