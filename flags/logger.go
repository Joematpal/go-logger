package flags

import (
	"strings"

	"github.com/joematpal/go-logger"
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
		Usage:   "enables the stacktrace after an error or higher log",
		Aliases: []string{"lst"},
		EnvVars: flagNamesToEnv(LogStacktrace, "LST"),
	},
	&cli.GenericFlag{
		Name:    LogEnv,
		Usage:   "values: prod, dev",
		Value:   logger.NewLogEnvEnum(),
		EnvVars: flagNamesToEnv(LogEnv),
	},
	&cli.GenericFlag{
		Name:    LogLevel,
		Usage:   "values: debug, info, warn, error, dpanic, panic, fatal",
		Value:   logger.NewLogLevelEnum(),
		EnvVars: flagNamesToEnv(LogLevel),
	},
	&cli.GenericFlag{
		Name:    LogEncoding,
		Usage:   "values: json, encoding",
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
