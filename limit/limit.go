package limit

import (
	"sync"
)

// Limit limit
type Limit struct {
	ch   chan struct{}
	wg   sync.WaitGroup
	once sync.Once
	Err  error
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

func (limit *Limit) Go(fc func() error) {
	limit.Add(1)
	go func() {
		defer limit.Done()
		if err := fc(); err != nil {
			limit.once.Do(func() {
				limit.Err = err
			})
		}
	}()
}

// Wait wait all groutine close
func (limit *Limit) Wait() error {
	limit.wg.Wait()
	return limit.Err
}

// Len len
func (limit *Limit) Len() int {
	return len(limit.ch)
}
