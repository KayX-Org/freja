package freja

import (
	"io"
)

//go:generate moq -out logger_mock_test.go . Logger
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Printf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warn(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

type DummyLogger struct{}

func (*DummyLogger) Debugf(format string, args ...interface{})   {}
func (*DummyLogger) Infof(format string, args ...interface{})    {}
func (*DummyLogger) Printf(format string, args ...interface{})   {}
func (*DummyLogger) Warnf(format string, args ...interface{})    {}
func (*DummyLogger) Warningf(format string, args ...interface{}) {}
func (*DummyLogger) Errorf(format string, args ...interface{})   {}
func (*DummyLogger) Fatalf(format string, args ...interface{})   {}
func (*DummyLogger) Panicf(format string, args ...interface{})   {}
func (*DummyLogger) Debug(args ...interface{})                   {}
func (*DummyLogger) Info(args ...interface{})                    {}
func (*DummyLogger) Print(args ...interface{})                   {}
func (*DummyLogger) Warn(args ...interface{})                    {}
func (*DummyLogger) Warning(args ...interface{})                 {}
func (*DummyLogger) Error(args ...interface{})                   {}
func (*DummyLogger) Fatal(args ...interface{})                   {}
func (*DummyLogger) Panic(args ...interface{})                   {}

type LogWriter interface {
	io.Writer
}

type LogWrite struct {
	log   Logger
	level string
}

func NewLogWrite(log Logger, level string) *LogWrite {
	return &LogWrite{
		log:   log,
		level: level,
	}
}

func (l *LogWrite) Write(p []byte) (int, error) {
	switch l.level {
	case "debug":
		l.log.Debug(string(p))
	case "info":
		l.log.Info(string(p))
	case "warn":
		l.log.Warn(string(p))
	case "error":
		l.log.Error(string(p))
	case "fatal":
		l.log.Fatalf(string(p))
	case "panic":
		l.log.Panic(string(p))
	default:
		l.log.Warn(string(p))
	}

	return len(p), nil
}
