package main

import (
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
			logr, err := logger.New()
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
				log.Fatal(err)
			}
			return nil
		},
	}
}

func main() {
	if err := NewApp().Run(os.Args); err != nil {
		panic(err)
	}
}
