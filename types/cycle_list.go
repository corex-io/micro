package types

import (
	"sync"
)

// CycleList CycleList
type CycleList struct {
	idx int
	max int
	arr []string
	cnt int
	mux sync.RWMutex
}

// NewCycleList create new cycleList
func NewCycleList(max int) *CycleList {
	return &CycleList{
		idx: 0,
		max: max,
		arr: make([]string, max, max),
	}

}

// Append one value
func (c *CycleList) Append(v string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.arr[c.idx] = v
	c.cnt++
	//c.idx = (c.idx + 1) % c.max
	if c.idx == c.max-1 {
		c.idx = 0
	} else {
		c.idx++
	}
}

// List dump List
func (c *CycleList) List() []string {
	c.mux.RLock()
	defer c.mux.RUnlock()
	v := make([]string, 0, c.max)

	if c.cnt > c.max-1 {
		v = append(v, c.arr[c.idx:]...)

	}
	v = append(v, c.arr[0:c.idx]...)
	return v
}
