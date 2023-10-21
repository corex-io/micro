package micro

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/corex-io/micro/log"

	"golang.org/x/sync/errgroup"
)

// service
type service struct {
	options Options
	once    sync.Once

	ctx   context.Context
	cause context.CancelCauseFunc

	init     []Runnable
	services []Runnable
}

// newService
func newService(opts ...Option) *service {
	svc := &service{
		options: newOptions(opts...),
	}
	return svc
}

// Prepare Init init
func (s *service) Prepare(opts ...Option) {
	for _, o := range opts {
		o(&s.options)
	}

	s.once.Do(func() {
	})
}

// Name xx
func (s *service) Name() string {
	return s.options.AppName
}

// Options options
func (s *service) Options() Options {
	return s.options
}

// Run run
func (s *service) Run() error {
	err := s.run()
	log.Infof("shut down, cost=%s, err=%v", time.Since(s.options.Uptime), err)
	return err
}

// Run run
func (s *service) run() error {
	log.Infof("name=%s, pid=%d, ip=%s", s.options.AppName, s.options.PID, s.options.LocIP)

	var group *errgroup.Group
	s.ctx, s.cause = context.WithCancelCause(context.Background())
	group, s.ctx = errgroup.WithContext(s.ctx)
	go s.notify()

	// if s.init != nil {
	// 	if err := s.init.Init(s.ctx); err != nil {
	// 		return err
	// 	}
	// }

	defer func() {
		for i := len(s.services); i > 0; i-- {
			if v, ok := s.services[i-1].(io.Closer); ok {
				v.Close()
			}
		}

	}()

	for _, svc := range s.services {
		if v, ok := svc.(Init); ok {
			if err := v.Init(s.ctx); err != nil {
				return err
			}
			log.Infof("> %T[%s] InitOK", svc, svc)
		}
	}

	for _, svc := range s.services {
		svc := svc
		group.Go(func() error {
			log.Infof("> %T[%s] Running...", svc, svc)
			defer func() {
				if err := recover(); err != nil {
					for _, txt := range bytes.SplitN(debug.Stack(), []byte("\n"), -1) {
						log.Errorf("%T: %s", svc, txt)
					}
					log.Errorf("%T:[%s] Exit. err=%v", svc, svc, err)
				} else {
					log.Infof("> %T[%s] Exit.", svc, svc)
				}
			}()
			return svc.Run(s.ctx)
		})
	}
	return group.Wait()
}

// notify
func (s *service) notify() {
	ignore := make(chan os.Signal, 1)
	sign := make(chan os.Signal, 1)

	signal.Notify(ignore, syscall.SIGHUP) // 终端挂起或者控制进程终止(hang up)
	signal.Notify(sign, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-s.ctx.Done():
		s.cause(s.ctx.Err())
	case sig := <-sign:
		s.cause(fmt.Errorf("got sign: %s", sig))
	}
	log.Warnf("waiting all grountine closed")
}

// Regist regist
func (s *service) Regist(runnable Runnable) {
	s.services = append(s.services, runnable)
}

// RegistFunc registFunc
func (s *service) RegistFunc(runFunc RunFunc) {
	s.services = append(s.services, runFunc)
}
