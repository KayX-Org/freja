package component

import (
	log "github.com/sirupsen/logrus"
)

type Log struct {
	log log.Logger
}

func NewLogger() *log.Entry {
	log.SetFormatter(&log.JSONFormatter{})
	serviceName := getEnv("SERVICE_NAME", "service")
	standardFields := log.Fields{
		"serviceName": serviceName,
	}

	return log.WithFields(standardFields)
}
