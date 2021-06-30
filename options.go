package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option interface {
	applyOption(*zap.Config) error
}

type applyOptionFunc func(*zap.Config) error

func (f applyOptionFunc) applyOption(o *zap.Config) error {
	return f(o)
}

func WithLevel(level string) Option {
	return applyOptionFunc(func(c *zap.Config) error {
		val, ok := LogLevelEnum_values[level]
		if !ok {
			return fmt.Errorf("invalid log level: %s", level)
		}
		c.Level = zap.NewAtomicLevelAt(zapcore.Level(val))
		return nil
	})
}

func WithEnv(env string) Option {
	return applyOptionFunc(func(c *zap.Config) error {
		if env == dev {
			c.Development = true
		}
		return nil
	})
}

func WithEncoding(encoding string) Option {
	return applyOptionFunc(func(c *zap.Config) error {
		c.Encoding = encoding
		return nil
	})
}

func WithLogStacktrace(lst bool) Option {
	return applyOptionFunc(func(c *zap.Config) error {
		c.DisableStacktrace = !lst
		return nil
	})
}
