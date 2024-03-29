package logger

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/joematpal/go-logger/event"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	writers []io.Writer
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

func writeByNewLine(factoryError FactoryError, reader io.Reader, writers ...io.Writer) error {
	return writeByNewLineWithContext(context.Background(), factoryError, reader, writers...)
}

func writeByNewLineWithContext(ctx context.Context, factoryError FactoryError, reader io.Reader, writers ...io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)

	eg, ctx := errgroup.WithContext(ctx)

	e := event.New()
	for _, writer := range writers {
		writer := writer
		eg.Go(func() error {
			itr := e.Subscribe()
			for {
				select {
				case data := <-itr.Data():
					if _, err := writer.Write(data); err != nil {
						cancel()
						return err
					}
				case <-ctx.Done():
					cancel()
					err := ctx.Err()
					if !errors.Is(err, context.Canceled) {
						return err
					}
					return nil
				}
			}
		})
	}

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		b := scanner.Bytes()
		b = append(b, '\n')
		e.Publish(b)
	}

	serr := scanner.Err()
	cancel()
	err := eg.Wait()
	if err != nil {
		return fmt.Errorf("errgroup: %w", err)
	}
	if serr != nil {
		serr = fmt.Errorf("scanner: %v", serr)
	}
	if err != nil || serr != nil {
		return fmt.Errorf("%v: %v", err, serr)
	}
	return nil
}

func writeByNewLineSync(factoryError FactoryError, reader io.Reader, writers ...io.Writer) error {
	wr := io.MultiWriter(writers...)
	if _, err := io.Copy(wr, reader); err != nil {
		factoryError.Error(err)
		return err
	}
	return nil
}
