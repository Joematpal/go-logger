package logger

import "go.uber.org/zap"

type Option interface {
	applyOption(*zap.Config) error
}

type applyOptionFunc func(*zap.Config) error

func (f applyOptionFunc) applyOption(o *zap.Config) error {
	return f(o)
}

func WithLevel(level LogLevel) Option {
	return applyOptionFunc(func(c *zap.Config) error {
		c.Level = zap.NewAtomicLevelAt(level)
		return nil
	})
}

func WithEnv(env LogEnv) Option {
	return applyOptionFunc(func(c *zap.Config) error {
		if env == Development {
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
		c.DisableStacktrace = lst
		return nil
	})
}
