package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	correlationID = "correlationId"
)

// Logger represent common interface for logging function
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	DPanic(args ...interface{})
	DPanicf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
}

type CorrelationLogger interface {
	Logger
	WithCorrelationID(id string) CorrelationLogger
}

type SugaredLogger interface {
	Logger
	Desugar() *zap.Logger
}

type logger struct {
	log           SugaredLogger
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

	logr, err := config.Build(
		zap.WithCaller(false),
	)
	if err != nil {
		return nil, fmt.Errorf("build: %v", err)
	}
	sugar := logr.Sugar()

	return &logger{
		log: sugar,
	}, err
}

func NewCorrelationLogger(opts ...Option) (CorrelationLogger, error) {
	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	for _, opt := range opts {
		if err := opt.applyOption(&config); err != nil {
			return nil, err
		}
	}

	logr, err := config.Build(
		zap.WithCaller(false),
	)
	if err != nil {
		return nil, fmt.Errorf("build: %v", err)
	}
	sugar := logr.Sugar()

	return &logger{
		log: sugar,
	}, err
}

func (l *logger) Debug(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Debug(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.Debug(argsToString(args))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Debug(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.Debugf(format, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().DPanic(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.DPanic(argsToString(args))
}

func (l *logger) DPanicf(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().DPanic(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.DPanicf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Error(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.Error(argsToString(args))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Error(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.Errorf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Fatal(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.Fatal(argsToString(args))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Fatal(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.Fatalf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Info(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.Info(argsToString(args))
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Info(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.Infof(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Warn(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.Warn(argsToString(args))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Warn(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.Warnf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Panic(argsToString(args), zap.String(correlationID, l.correlationID))
	}
	l.log.Panic(argsToString(args))
}

func (l *logger) Panicf(format string, args ...interface{}) {
	if l.correlationID != "" {
		l.log.Desugar().Panic(fmt.Sprintf(format, args...), zap.String(correlationID, l.correlationID))
	}
	l.log.Panicf(format, args...)
}

func (l *logger) WithCorrelationID(id string) CorrelationLogger {
	return &logger{
		log:           l.log,
		correlationID: id,
	}
}

func argsToString(args []interface{}) string {
	var sb strings.Builder
	for i, arg := range args {

		if len(args)-1 == i {
			sb.WriteString(fmt.Sprintf("%v", arg))
		} else {
			sb.WriteString(fmt.Sprintf("%v ", arg))
		}
	}

	return sb.String()
}
