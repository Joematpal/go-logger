package logger

import (
	"fmt"
	"strings"
)

/*
LOG LEVEL
*/

const (
	debug    = "debug"
	info     = "info"
	warn     = "warn"
	errorStr = "error"
	dpanic   = "dpanic"
	panic    = "panic"
	fatal    = "fatal"
)

type LogLevelEnum struct {
	Enum     []string
	Default  int
	selected int
}

func NewLogLevelEnum() *LogLevelEnum {
	return &LogLevelEnum{
		Enum: []string{
			debug,
			info,
			warn,
			errorStr,
			dpanic,
			panic,
			fatal,
		},
	}
}

var LogLevelEnum_values = map[string]int{
	debug:    int(DebugLevel),
	info:     int(InfoLevel),
	warn:     int(WarnLevel),
	errorStr: int(ErrorLevel),
	dpanic:   int(DPanicLevel),
	panic:    int(PanicLevel),
	fatal:    int(FatalLevel),
}

var LogLevelEnum_keys = map[int]string{
	int(DebugLevel):  debug,
	int(InfoLevel):   info,
	int(WarnLevel):   warn,
	int(ErrorLevel):  errorStr,
	int(DPanicLevel): dpanic,
	int(PanicLevel):  panic,
	int(FatalLevel):  fatal,
}

func (e *LogLevelEnum) Set(value string) error {
	if val, ok := LogLevelEnum_values[value]; ok {
		e.selected = val
		return nil
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e *LogLevelEnum) String() string {
	if val, ok := LogLevelEnum_keys[e.selected]; ok {
		return val
	}
	return ""
}

/*
ENVIRONMENT
*/

type LogEnv int

const (
	Development LogEnv = iota
	Production
	prod = "prod"
	dev  = "dev"
)

func (e LogEnv) String() string {
	switch e {
	case Production:
		return prod
	case Development:
		return dev
	}
	return ""
}

func ParseLogEnv(env string) LogEnv {
	switch strings.ToLower(env) {
	case prod:
		return Production
	case dev:
		return Development
	}
	return Development
}

type LogEnvEnum struct {
	Enum     []string
	Default  string
	selected string
}

func NewLogEnvEnum() *LogEnvEnum {
	return &LogEnvEnum{
		Enum: []string{
			prod,
			dev,
		},
		Default: dev,
	}
}

var LogEnvEnum_values = map[string]LogEnv{
	prod: Production,
	dev:  Development,
}

var LogEnvEnum_keys = map[LogEnv]string{
	Production:  prod,
	Development: dev,
}

func (e *LogEnvEnum) Set(value string) error {
	if val, ok := LogEnvEnum_values[strings.ToLower(value)]; ok {
		e.selected = val.String()
		return nil
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e *LogEnvEnum) String() string {
	if val, ok := LogEnvEnum_keys[ParseLogEnv(e.selected)]; ok {
		return val
	}
	return ""
}

/*
ENCODING
*/

type LogEncoding int

const (
	JSON LogEncoding = iota
	Console
	json    = "json"
	console = "console"
)

func (e LogEncoding) String() string {
	switch e {
	case JSON:
		return json
	case Console:
		return console
	}
	return ""
}

func ParseLogEncoding(env string) LogEncoding {
	switch strings.ToLower(env) {
	case json:
		return JSON
	case console:
		return Console
	}
	return Console
}

type LogEncodingEnum struct {
	Enum     []string
	Default  string
	selected string
}

func NewLogEncodingEnum() *LogEncodingEnum {
	return &LogEncodingEnum{
		Enum: []string{
			json,
			console,
		},
		Default: console,
	}
}

var LogEncodingEnum_values = map[string]LogEncoding{
	json:    JSON,
	console: Console,
}

var LogEncodingEnum_keys = map[LogEncoding]string{
	JSON:    json,
	Console: console,
}

func (e *LogEncodingEnum) Set(value string) error {
	if val, ok := LogEncodingEnum_values[strings.ToLower(value)]; ok {
		e.selected = val.String()
		return nil
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e *LogEncodingEnum) String() string {
	if val, ok := LogEncodingEnum_keys[ParseLogEncoding(e.selected)]; ok {
		return val
	}
	return ""
}
