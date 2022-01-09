package limit

import (
	"context"
	"errors"
	"time"

	"LargeScale/micro/log"
)

var StopErr = errors.New("stop limit.Every")

// Every every
type Every struct {
	duration time.Duration
	f func(ctx context.Context) error
	stop     chan error
}

// NewEvery newEvery
func NewEvery(duration time.Duration, f func(ctx context.Context) error) *Every {
	return &Every{
		duration: duration,
		f: f,
		stop:     make(chan error),
	}
}
// Stop stop
func (v *Every) Stop(err error) {
	v.stop <- err
}

// Run 间隔every秒开始下一次
func (v *Every) Run(ctx context.Context) error {
	t := time.NewTimer(0 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := v.f(ctx); err != nil {
				if errors.Is(err, StopErr) {
					return err
				}
				log.Errorf("%s", err)
			}
			t.Reset(v.duration)
		case <-ctx.Done():
			return ctx.Err()
		case err :=<-v.stop:
			return err
		}
	}
}

// Run 每隔every秒开始执行下一次，如果一次任务的运行时长超过了every，那么下一次直接发起
func (v *Every) Sched(ctx context.Context) error {
	t := time.NewTimer(0 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			t.Reset(v.duration)
			if err := v.f(ctx); err != nil {
				if errors.Is(err, StopErr) {
					return err
				}
				log.Errorf("%s", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		case err :=<-v.stop:
			return err
		}
	}
}

func (v Every) Name() string{
	return "Every"
}