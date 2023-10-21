package micro

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"
)

func Pprof(ctx context.Context, duration time.Duration) error {
	if duration <= 5*time.Second || duration >= 10*time.Minute {
		duration = 30 * time.Second
	}
	f, err := os.OpenFile("pprof.out", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "start cpu profile: %s\n", time.Now().Format(time.RFC3339))
	defer pprof.StopCPUProfile()
	defer func() {
		fmt.Fprintf(os.Stdout, "Please use `go tool -http:0.0.0.0:80 pprof prof.out`")
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}

func Trace(ctx context.Context, duration time.Duration) error {
	if duration <= 5*time.Second || duration >= 10*time.Minute {
		duration = 30 * time.Second
	}
	f, err := os.OpenFile("trace.out", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	trace.Start(f)
	defer trace.Stop()
	defer func() {
		fmt.Fprintf(os.Stdout, "Please use `go tool -http:0.0.0.0:80 trace trace.out`")
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(duration):
		return nil
	}
}
