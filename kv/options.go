package kv

// Options options
type Options struct {
	bucket string
}

func newOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Option option
type Option func(*Options)

// Bucket bucket
func Bucket(bucket string) Option {
	return func(opt *Options) {
		opt.bucket = bucket
	}
}
