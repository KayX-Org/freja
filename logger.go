package freja

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
