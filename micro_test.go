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
	service := New(Name("test"))
	service.RegistFunc(func(ctx context.Context) error {
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
	if err := service.Run(); err != nil {
		t.Error(err)
	}
}
