package logger

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap/zapcore"
)

var (
	errNoEncoderNameSpecified = errors.New("no encoder name specified")

	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) (zapcore.Encoder, error){
		"console": func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return zapcore.NewConsoleEncoder(encoderConfig), nil
		},
		"json": func(encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
			return zapcore.NewJSONEncoder(encoderConfig), nil
		},
	}
	_encoderMutex sync.RWMutex
)

func newEncoder(name string, encoderConfig zapcore.EncoderConfig) (zapcore.Encoder, error) {
	if encoderConfig.TimeKey != "" && encoderConfig.EncodeTime == nil {
		return nil, fmt.Errorf("missing EncodeTime in EncoderConfig")
	}

	_encoderMutex.RLock()
	defer _encoderMutex.RUnlock()
	if name == "" {
		return nil, errNoEncoderNameSpecified
	}
	constructor, ok := _encoderNameToConstructor[name]
	if !ok {
		return nil, fmt.Errorf("no encoder registered for name %q", name)
	}
	return constructor(encoderConfig)
}
