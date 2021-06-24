package logger

type Options struct {
}

type Option interface {
	applyOption(*Options) error
}

type applyOptionFunc func(*Options) error

func (f applyOptionFunc) applyOption(o *Options) error {
	return f(o)
}
