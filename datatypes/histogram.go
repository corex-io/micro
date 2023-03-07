package datatypes

import "sync/atomic"

type Histogram struct {
	Le   []float64
	hist []int64
}

func NewHistogram(le ...float64) *Histogram {
	return &Histogram{
		Le:   le,
		hist: make([]int64, len(le), len(le)),
	}
}

func (his *Histogram) Compute(value float64) {
	for i, v := range his.Le {
		if value > v {
			continue
		}
		atomic.AddInt64(&his.hist[i], 1)
		return
	}
}
