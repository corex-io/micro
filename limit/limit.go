package limit

import (
	"fmt"
	"sync"

	"github.com/corex-io/micro/log"
)

// Limit limit
type Limit struct {
	ch chan struct{}
	wg sync.WaitGroup
}

// NewLimit newLimit
func NewLimit(max int) *Limit {
	return &Limit{
		ch: make(chan struct{}, max),
	}
}

func (limit *Limit) Add(delta int) {
	limit.ch <- struct{}{}
	limit.wg.Add(1)
}

func (limit *Limit) Done() {
	<-limit.ch
	limit.wg.Done()
}

func (limit *Limit) Go(fc func() error) error {
	limit.Add(1)
	go func() error {
		defer limit.Done()
		return fc()

	}()
	return nil
}

// execute run
func (limit *Limit) try(fc func() error) error {
	select {
	case limit.ch <- struct{}{}:
		limit.wg.Add(1)
		defer func() {
			<-limit.ch
			limit.wg.Done()
		}()
		return fc()
	default:
		return fmt.Errorf("limited[%d]", cap(limit.ch))
	}
}

// TryGo try run once, if fail exit
func (limit *Limit) TryGo(key string, fc func() error) {
	go func() {
		if err := limit.try(fc); err != nil {
			log.Errorf("key=%s, %v", key, err.Error())
		}
	}()
}

// Wait wait all groutine close
func (limit *Limit) Wait() {
	limit.wg.Wait()
}

// Len len
func (limit *Limit) Len() int {
	return len(limit.ch)
}
