package datatypes

import (
	"sort"
	"sync"
)

type Bucket struct {
	vList []float64
	once  sync.Once
}

func NewBucket2(v ...*float64) *Bucket {
	bucket := Bucket{}
	for _, x := range v {
		if x != nil {
			bucket.vList = append(bucket.vList, *x)
		}
	}
	return &bucket
}

func NewBucket(v ...float64) *Bucket {
	return &Bucket{
		vList: v,
	}
}

func (b *Bucket) Add(v ...float64) {
	b.vList = append(b.vList, v...)
}

func (b *Bucket) Less(i, j int) bool {
	return b.vList[i] < b.vList[j]
}

func (b *Bucket) Len() int {
	return len(b.vList)
}

func (b *Bucket) Swap(i, j int) {
	b.vList[i], b.vList[j] = b.vList[j], b.vList[i]
}

func (b *Bucket) SUM() float64 {
	sum := 0.0
	for _, i := range b.vList {
		sum += i
	}
	return sum
}

// AVG AVG
func (b *Bucket) AVG() float64 {
	return b.SUM() / float64(b.Len())
}

// MIN MIN
func (b *Bucket) MIN() float64 {
	b.once.Do(func() { sort.Sort(b) })
	return b.vList[0]
}

// MAX MAX
func (b *Bucket) MAX() float64 {
	b.once.Do(func() { sort.Sort(b) })
	return b.vList[b.Len()-1]
}

// MED 中位数:Median
func (b *Bucket) MED() float64 {
	b.once.Do(func() { sort.Sort(b) })
	return b.vList[b.Len()/2]
}

func (b *Bucket) P(p float64) float64 {
	b.once.Do(func() { sort.Sort(b) })
	i := int(p * float64(b.Len()))
	return b.vList[i]
}

func (b *Bucket) P99() float64 {
	return b.P(0.99)
}

func (b *Bucket) P95() float64 {
	return b.P(0.95)
}

func (b *Bucket) P90() float64 {
	return b.P(0.90)
}
