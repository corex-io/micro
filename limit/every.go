package limit

import (
	"context"
	"errors"
	"time"

	"golang.org/x/sync/errgroup"
)

var StopErr = errors.New("stop limit.Every")

// Every every
type Every struct {
	wait time.Duration
	f    func(ctx context.Context) error
	st   func(ctx context.Context) error
}

// NewEvery newEvery
func NewEvery(wait time.Duration, f func(ctx context.Context) error) *Every {
	return &Every{
		wait: wait,
		f:    f,
	}
}

func (v *Every) Sentry(f func(ctx context.Context) error) {
	v.st = f
}

// Run 每隔every秒开始执行下一次，如果一次任务的运行时长超过了every，那么下一次直接发起
func (v *Every) Sched(ctx context.Context, exit func(error) bool) error {
	group, ctx := errgroup.WithContext(ctx)
	defer group.Wait()
	if v.st != nil {
		group.Go(func() error {
			return v.st(ctx)
		})
	}

	t := time.NewTimer(0 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := v.f(ctx); exit(err) {
				return err
			}
			t.Reset(v.wait)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (v Every) Name() string {
	return "Every"
}
