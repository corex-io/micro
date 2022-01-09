package log

// Options opts
type Options struct {
	timeFormat string
	msgFormat  string
}

// Option func
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		timeFormat: "2006/01/02 15:04:05.000",
	}
	return withOptions(opt, opts...)
}

// WithOptions withOptions
func withOptions(opt Options, opts ...Option) Options {
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Format logformat
func Format(format string) Option {
	return func(o *Options) {
		o.msgFormat = format
	}
}

// TimeFormat timeFormat
func TimeFormat(timeFormat string) Option {
	return func(o *Options) {
		o.timeFormat = timeFormat
	}
}
