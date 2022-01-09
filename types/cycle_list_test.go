package types

import (
	"fmt"
	"testing"
)

func BenchmarkCycleList(b *testing.B) {
	list := NewCycleList(10000)
	for i := 0; i < b.N; i++ {
		list.Append("1")
	}
}

func TestCycleList(t *testing.T) {
	list := NewCycleList(5)
	for i := 0; i < 20; i++ {
		list.Append(fmt.Sprintf("%d", i))
		fmt.Println(list.List())
	}
}
