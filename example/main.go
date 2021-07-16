package main

import (
	"fmt"
	"log"
	"os"

	logger "github.com/joematpal/go-logger"
	nested "github.com/joematpal/go-logger/example/nested_package"
	"github.com/joematpal/go-logger/flags"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Flags: flags.LogFlags,
		Action: func(c *cli.Context) error {
			opts := []logger.Option{
				logger.WithEnv(c.String(flags.LogEnv)),
				logger.WithLevel(c.String(flags.LogLevel)),
				logger.WithLogStacktrace(c.Bool(flags.LogStacktrace)),
				logger.WithEncoding(c.String(flags.LogEncoding)),
			}

			if encoding := c.String(flags.LogEncoding); encoding != "" {
				opts = append(opts, logger.WithEncoding(encoding))
			}

			logr, err := logger.NewCorrelationLogger(opts...)
			if err != nil {
				log.Fatal(err)
			}
			logr.Debug("test", "more")
			logr.Debug("1", "2 ")
			logr.Debugf("key=%s", "value")

			// With Correlation ID
			clogr := logr.WithCorrelationID("test_id")
			clogr.Debug("debug stuff", "and another message")
			clogr.Debugf("magic=%s", "spell")
			clogr.Info("test")
			clogr.Infof("somthing %s", "special")
			logr.Info("what")

			logr.Error("test")
			logr.Info("after test")

			client := &nested.Client{
				Logger: logr,
			}

			if err := client.Something(); err != nil {
				return fmt.Errorf("client.Something: %v", err)
			}

			return nil
		},
	}
}

func main() {
	if err := NewApp().Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
