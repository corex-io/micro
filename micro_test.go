package micro

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/corex-io/micro/datatypes"
	"github.com/corex-io/micro/log"
)

func TestService(t *testing.T) {
	app := New(Name("test"))
	app.RegistFunc(func(ctx context.Context) error {
		fmt.Println("init...")
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				log.Infof("%s", datatypes.Now())
				time.Sleep(1 * time.Second)
			}
		}
	})
	app.RegistFunc(func(ctx context.Context) error {
		return Pprof(ctx, 10*time.Minute)
	})
	app.RegistFunc(func(ctx context.Context) error {
		return Trace(ctx, 10*time.Minute)
	})
	if err := app.Run(); err != nil {
		t.Error(err)
	}
}
