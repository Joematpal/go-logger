package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	correlationID = "correlation_id"
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
	fields        fields
}

type FieldLogger interface {
	// Logger
	// WithCorrelationID(id string) CorrelationLogger
	CorrelationLogger
	WithField(key string, value interface{}) FieldLogger
	WithFields(in ...Field) FieldLogger
}

func New(opts ...Option) (Logger, error) {
	zapc := zap.NewProductionConfig()
	config := &Config{
		zap: &zapc,
	}
	config.zap.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	for _, opt := range opts {
		if err := opt.applyOption(config); err != nil {
			return nil, err
		}
	}

	buildOpts := []zap.Option{
		zap.WithCaller(false),
	}
	if config.writer != nil {
		f, err := newCore(config)
		if err != nil {
			return nil, err
		}
		buildOpts = append(buildOpts, zap.WrapCore(f))
	}
	logr, err := config.zap.Build(
		buildOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("build: %v", err)
	}
	sugar := logr.Sugar()

	return &logger{
		log:    sugar,
		fields: fields{},
	}, err
}

func NewCorrelationLogger(opts ...Option) (CorrelationLogger, error) {
	zapc := zap.NewProductionConfig()
	config := &Config{
		zap: &zapc,
	}
	config.zap.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	for _, opt := range opts {
		if err := opt.applyOption(config); err != nil {
			return nil, err
		}
	}

	buildOpts := []zap.Option{
		zap.WithCaller(false),
	}
	if config.writer != nil {
		f, err := newCore(config)
		if err != nil {
			return nil, err
		}
		buildOpts = append(buildOpts, zap.WrapCore(f))
	}

	logr, err := config.zap.Build(
		buildOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("build: %v", err)
	}
	sugar := logr.Sugar()

	return &logger{
		log:    sugar,
		fields: fields{},
	}, err
}

func (l *logger) Debug(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Debug(argsToString(args), fields...)
		return
	}
	l.log.Debug(argsToString(args))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Debug(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.Debugf(format, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().DPanic(argsToString(args), fields...)
		return
	}
	l.log.DPanic(argsToString(args))
}

func (l *logger) DPanicf(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().DPanic(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.DPanicf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Error(argsToString(args), fields...)
		return
	}
	l.log.Error(argsToString(args))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Error(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.Errorf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Fatal(argsToString(args), fields...)
		return
	}
	l.log.Fatal(argsToString(args))
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Fatal(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.Fatalf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.
			Desugar().
			Info(argsToString(args), fields...)
		return
	}
	l.log.Info(argsToString(args))
}

func (l *logger) Infof(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Info(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.Infof(format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Warn(argsToString(args), fields...)
		return
	}
	l.log.Warn(argsToString(args))
}

func (l *logger) Warnf(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Warn(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.Warnf(format, args...)
}

func (l *logger) Panic(args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)

	if len(fields) > 0 {
		l.log.Desugar().Panic(argsToString(args), fields...)
		return
	}
	l.log.Panic(argsToString(args))
}

func (l *logger) Panicf(format string, args ...interface{}) {
	fields := getFields(l.correlationID, l.fields)
	if len(fields) > 0 {
		l.log.Desugar().Panic(fmt.Sprintf(format, args...), fields...)
		return
	}
	l.log.Panicf(format, args...)
}

func (l *logger) WithCorrelationID(id string) CorrelationLogger {
	return &logger{
		log:           l.log,
		correlationID: id,
	}
}

func getFields(cID string, fields fields) []zapcore.Field {
	out := []zapcore.Field{}
	if cID != "" {
		out = append(out, zap.String(correlationID, cID))
	}
	for _, field := range fields {
		if field.key == correlationID {
			continue
		}
		out = append(out, zap.Any(field.key, field.value))
	}
	return out
}

// Fields
type field struct {
	key   string
	value interface{}
}

type Field interface {
	Key() string
	Value() interface{}
}

type fields = []field

func (l *logger) WithField(key string, value interface{}) FieldLogger {
	fields := l.fields
	fields = append(fields, field{key, value})

	return &logger{
		log:           l.log,
		correlationID: l.correlationID,
		fields:        fields,
	}
}

func (l *logger) WithFields(in ...Field) FieldLogger {
	fields := l.fields
	for _, f := range in {
		fields = append(fields, field{f.Key(), f.Value()})
	}

	return &logger{
		log:           l.log,
		correlationID: l.correlationID,
		fields:        fields,
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

type newCoreFunc = func(c zapcore.Core) zapcore.Core

func newCore(config *Config) (newCoreFunc, error) {
	enc, err := newEncoder(config.zap.Encoding, config.zap.EncoderConfig)
	if err != nil {
		return nil, err
	}
	return func(c zapcore.Core) zapcore.Core {
		return zapcore.NewCore(enc, zapcore.AddSync(config.writer), zapcore.DebugLevel)
	}, err
}
