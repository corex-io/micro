package limit

import (
	"sync/atomic"
)

// Counter 计数器接口
type Counter interface {
	Add(int64) int64
	Count() int64
	Reset()
}

// LoadCounter 基数计数器
func LoadCounter(base int64) Counter {
	return &count{c: base}
}

// NewCounter 计数器
func NewCounter() Counter {
	return &count{}
}

type count struct {
	c int64
}

// Add count
func (c *count) Add(delta int64) int64 {
	return atomic.AddInt64(&c.c, delta)
}

// Count counter
func (c *count) Count() int64 {
	return atomic.LoadInt64(&c.c)
}

// Reset counter
func (c *count) Reset() {
	atomic.StoreInt64(&c.c, 0)
}
