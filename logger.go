package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Logger represent common interface for logging function
type Logger interface {
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatal(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Debug(args ...interface{})
	WithCorrelationID(id string) Logger
}

type logger struct {
	log           *zap.SugaredLogger
	correlationID string
}

func New() (Logger, error) {
	logr, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := logr.Sugar()

	return &logger{
		log: sugar,
	}, err
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

func (l *logger) WithCorrelationID(id string) Logger {
	return &logger{
		log:           l.log,
		correlationID: id,
	}
}
