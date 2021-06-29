package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represent common interface for logging function
type Logger interface {
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	DPanicf(format string, args ...interface{})
	DPanic(args ...interface{})
	Panicf(format string, args ...interface{})
	Panic(args ...interface{})
	Warnf(format string, args ...interface{})

	WithCorrelationID(id string) Logger
}

type logger struct {
	log           *zap.SugaredLogger
	correlationID string
}

func New(opts ...Option) (Logger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	for _, opt := range opts {
		if err := opt.applyOption(&config); err != nil {
			return nil, err
		}
	}

	logr, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("build: %v", err)
	}
	sugar := logr.Sugar()

	return &logger{
		log: sugar,
	}, err
}

func (l *logger) Error(args ...interface{}) {
	if l.correlationID != "" {
		args = append([]interface{}{fmt.Sprintf("correlationID=%s", l.correlationID)}, args...)
	}
	l.log.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.Fatalf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	if l.correlationID != "" {
		args = append([]interface{}{fmt.Sprintf("correlationID=%s", l.correlationID)}, args...)
	}
	l.log.Fatal(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.Infof(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	if l.correlationID != "" {
		args = append([]interface{}{fmt.Sprintf("correlationID=%s", l.correlationID)}, args...)
	}
	l.log.Info(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.Warnf(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.Debugf(format, args...)
}

func (l *logger) Debug(args ...interface{}) {
	if l.correlationID != "" {
		args = append([]interface{}{fmt.Sprintf("correlationID=%s", l.correlationID)}, args...)
	}
	l.log.Debug(args...)
}

func (l *logger) DPanicf(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.DPanicf(format, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	if l.correlationID != "" {
		args = append([]interface{}{fmt.Sprintf("correlationID=%s", l.correlationID)}, args...)
	}
	l.log.DPanic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	if l.correlationID != "" {
		format = fmt.Sprintf("correlationID=%s %s", l.correlationID, format)
	}
	l.log.Panicf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	if l.correlationID != "" {
		args = append([]interface{}{fmt.Sprintf("correlationID=%s", l.correlationID)}, args...)
	}
	l.log.Panic(args...)
}

func (l *logger) WithCorrelationID(id string) Logger {
	return &logger{
		log:           l.log,
		correlationID: id,
	}
}
