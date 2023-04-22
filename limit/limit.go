package limit

import (
	"sync"
)

// Limit limit
type Limit struct {
	ch   chan struct{}
	wg   sync.WaitGroup
	once sync.Once
	Err  chan error
}

// NewLimit newLimit
func NewLimit(max int) *Limit {
	return &Limit{
		ch:  make(chan struct{}, max),
		Err: make(chan error, max),
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

func (limit *Limit) Go(fc func() error) {
	limit.Add(1)
	go func() {
		defer limit.Done()
		if err := fc(); err != nil {
			select {
			case limit.Err <- err:
			default:
			}
		}
	}()
}

// Wait wait all groutine close
func (limit *Limit) Wait() error {
	limit.wg.Wait()
	select {
	case err := <-limit.Err:
		return err
	default:
		return nil
	}
}

// Len len
func (limit *Limit) Len() int {
	return len(limit.ch)
}
