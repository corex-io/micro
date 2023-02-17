package limit

import (
	"context"
	"errors"
	"time"
)

var StopErr = errors.New("stop limit.Every")

// Every every
type Every struct {
	wait time.Duration
	fn   func(context.Context) error
	exit func(error) bool
}

// NewEvery newEvery
func NewEvery(wait time.Duration, fn func(ctx context.Context) error) *Every {
	return &Every{
		wait: wait,
		fn:   fn,
	}
}

func (v *Every) Exit(exit func(error) bool) *Every {
	v.exit = exit
	return v
}

// Run 每隔every秒开始执行下一次，如果一次任务的运行时长超过了every，那么下一次直接发起
func (v *Every) Run(ctx context.Context) error {
	t := time.NewTimer(0 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := v.fn(ctx); v.exit != nil && v.exit(err) {
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
