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
