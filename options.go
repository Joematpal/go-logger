package logger

import (
	"fmt"
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	writers []io.Writer
	fields  fields
	logFile string
	zap     *zap.Config
}

type Option interface {
	applyOption(*Config) error
}

type applyOptionFunc func(*Config) error

func (f applyOptionFunc) applyOption(c *Config) error {
	return f(c)
}

func WithLevel(level string) Option {
	return applyOptionFunc(func(c *Config) error {
		val, ok := LogLevelEnum_values[level]
		if !ok {
			return fmt.Errorf("invalid log level: %s", level)
		}
		c.zap.Level = zap.NewAtomicLevelAt(zapcore.Level(val))
		return nil
	})
}

func WithEnv(env string) Option {
	return applyOptionFunc(func(c *Config) error {
		if env == dev {
			c.zap.Development = true
		}
		return nil
	})
}

func WithEncoding(encoding string) Option {
	return applyOptionFunc(func(c *Config) error {
		c.zap.Encoding = encoding
		return nil
	})
}

func WithLogStacktrace(lst bool) Option {
	return applyOptionFunc(func(c *Config) error {
		c.zap.DisableStacktrace = !lst
		return nil
	})
}

func withWriter(writer io.Writer) Option {
	return applyOptionFunc(func(c *Config) error {
		c.writers = append(c.writers, writer)
		return nil
	})
}

func WithInitialFields(fields map[string]interface{}) Option {
	return applyOptionFunc(func(c *Config) error {
		c.zap.InitialFields = fields
		return nil
	})
}

func WithLogFile(logFile string) Option {
	return applyOptionFunc(func(c *Config) error {
		c.zap.OutputPaths = append(c.zap.OutputPaths, logFile)
		return nil
	})
}

func WithWriters(writers ...io.Writer) Option {
	return applyOptionFunc(func(c *Config) error {
		c.writers = append(c.writers, writers...)
		return nil
	})
}

func writeByNewLine(debugger Debugger, reader io.Reader, writers ...io.Writer) {
	go func() {
		for i, writer := range writers {
			if _, err := io.Copy(writer, reader); err != nil {
				debugger.Debugf("writeByNewLine: writer[%d]: %v", i, err)
			}
		}
	}()
}
