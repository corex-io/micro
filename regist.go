package micro

import "context"

// Runnable implements run method
type Runnable interface {
	Run(context.Context) error
}

// RunFunc implements Runnable interface
type RunFunc func(context.Context) error

// Run run
func (run RunFunc) Run(ctx context.Context) error {
	return run(ctx)
}

// Init interface
type Init interface {
	Init(context.Context) error
}

// InitFunc implements Runnable interface
type InitFunc func(context.Context) error

// Init init
func (init InitFunc) Init(ctx context.Context) error {
	return init(ctx)
}

type warpServiceFunc struct {
	name string
	f    func(context.Context) error
}

// String implement Stringer
func (s *warpServiceFunc) String() string {
	return s.name
}

// GoString implement GoStringer
func (s *warpServiceFunc) GoString() string {
	return s.name
}

// Run fc
func (s *warpServiceFunc) Run(ctx context.Context) error {
	return s.f(ctx)
}

// F warp func service
func F(name string, fc func(context.Context) error) *warpServiceFunc {
	return &warpServiceFunc{name: name, f: fc}
}
