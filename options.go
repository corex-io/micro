package micro

import (
	"context"
	"github.com/corex-io/micro/common"
	"os"
	"path"
	"time"
)

// Options options
type Options struct {
	AppName string // 程序名称
	PID     int    // 进程ID
	LocIP   string // 地IP
	Uptime  time.Time

	Context context.Context
}

// Option function
type Option func(*Options)

func newOptions(opts ...Option) Options {

	opt := Options{
		AppName: path.Base(os.Args[0]),
		PID:     os.Getpid(),
		Uptime:  time.Now(),
		LocIP:   common.IP,
		Context: context.Background(),
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Name service name
func Name(name string) Option {
	return func(opt *Options) {
		opt.AppName = name
	}
}

// Context specifies a context for the service.
func Context(ctx context.Context) Option {
	return func(opt *Options) {
		opt.Context = ctx
	}
}

// var mainFunction = func(ctx context.Context) error {
// 	<-ctx.Done()
// 	fmt.Println("main context cancel")
// 	return nil

// }

// // Regist regist
// func (e *Options) Regist(runnable Runnable) {
// 	e.services = append(e.services, runnable)
// }

// // RegistFunc registFunc
// func (e *Options) RegistFunc(runFunc RunFunc) {
// 	e.services = append(e.services, runFunc)
// }

// // SetPrepare set prepare
// func (e *Options) SetPrepare(prepare RunFunc) {
// 	e.prepare = prepare
// }

// // SetFunc set func
// func (e *Options) SetFunc(mainFunc RunFunc) {
// 	e.main = mainFunc
// }

// // SetExit set exit
// func (e *Options) SetExit(exitFunc RunFunc) {
// 	e.exit = exitFunc
// }

// // Run run
// func (e *Options) Run(ctx context.Context) {
// 	e.Log.Infof("%v", e.run(ctx))
// }

// func (e *Options) run(ctx context.Context) error {
// 	e.Log.Infof("name=%s, pid=%d, ip=%s", e.AppName, e.PID, e.LocIP)

// 	var wg sync.WaitGroup
// 	defer wg.Wait()

// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	if err := e.prepare(ctx); err != nil {
// 		e.Log.Errorf("Prepare: %v", err)
// 		return err
// 	}

// 	mainExit := make(chan error, 1)
// 	wg.Add(1)
// 	go func(mainExit chan<- error) {
// 		defer wg.Done()
// 		mainExit <- e.main(ctx)
// 	}(mainExit)

// 	defer e.exit(ctx) // nolint: errcheck

// 	for _, service := range e.services {
// 		wg.Add(1)
// 		go func(service Runnable) {
// 			defer wg.Done()
// 			err := service.Run(ctx)
// 			e.Log.Infof("%#v, %v", service.Run, err)
// 		}(service)
// 	}

// 	sign := make(chan os.Signal, 1)
// 	signal.Notify(sign, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
// 	signal.Ignore(syscall.SIGTERM)

// 	select {
// 	case s := <-sign:
// 		fmt.Println("cancel context")
// 		cancel()
// 		return fmt.Errorf("sign: %v", s)
// 	case err := <-mainExit:
// 		return err
// 	}
// }
