package flags

import (
	"strings"

	"github.com/digital-dream-labs/go-logger"
	cli "github.com/urfave/cli/v2"
)

var (
	LogEnv        = "log-env"
	LogLevel      = "log-level"
	LogStacktrace = "log-stacktrace"
	LogEncoding   = "log-encoding"
)

var LogFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    LogStacktrace,
		Value:   true,
		Aliases: []string{"lst"},
		EnvVars: flagNamesToEnv(LogStacktrace, "LST"),
	},
	&cli.GenericFlag{
		Name:    LogEnv,
		Value:   logger.NewLogEnvEnum(),
		EnvVars: flagNamesToEnv(LogEnv),
	},
	&cli.GenericFlag{
		Name:    LogLevel,
		Value:   logger.NewLogLevelEnum(),
		EnvVars: flagNamesToEnv(LogLevel),
	},
	&cli.GenericFlag{
		Name:    LogEncoding,
		Value:   logger.NewLogEncodingEnum(),
		EnvVars: flagNamesToEnv(LogEncoding),
	},
}

func flagNamesToEnv(names ...string) []string {
	out := []string{}
	for _, name := range names {
		out = append(out, flagNameToEnv(name))
	}
	return out
}

func flagNameToEnv(name string) string {
	return strings.ReplaceAll(strings.ToUpper(name), "-", "_")
}
