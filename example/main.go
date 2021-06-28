package main

import (
	"fmt"
	"log"
	"os"

	logger "github.com/digital-dream-labs/go-logger"
	nested "github.com/digital-dream-labs/go-logger/example/nested_package"
	"github.com/digital-dream-labs/go-logger/flags"
	"github.com/urfave/cli/v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Flags: flags.LogFlags,
		Action: func(c *cli.Context) error {
			opts := []logger.Option{
				logger.WithEnv(logger.LogEnv(c.Int(flags.LogEnv))),
				logger.WithLevel(logger.LogLevel(c.Int(flags.LogLevel))),
				logger.WithLogStacktrace(c.Bool(flags.LogStacktrace)),
			}

			if encoding := c.String(flags.LogEncoding); encoding != "" {
				opts = append(opts, logger.WithEncoding(encoding))
			}

			logr, err := logger.New(opts...)
			if err != nil {
				log.Fatal(err)
			}
			logr.Debug("test ", "more")
			logr.Debug("1 ", "2 ")
			logr.Debugf("key=%s", "value")

			// With Correlation ID

			clogr := logr.WithCorrelationID("test_id")
			clogr.Debug("debug stuff ", "and another message")
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
