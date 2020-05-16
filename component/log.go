package component

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

func NewLogger() *log.Entry {
	log.SetFormatter(&log.JSONFormatter{})
	serviceName := getEnv("SERVICE_NAME", "service")
	standardFields := log.Fields{
		"serviceName": serviceName,
	}
	switch strings.ToLower(getEnv("LOG_LEVEL", "warn")) {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.WarnLevel)
	}

	return log.WithFields(standardFields)
}
