package limit

import (
	"fmt"
	"testing"
	"time"
)

func TestCounter(t *testing.T) {
	c := NewCounter()
	for i := 0; i < 100; i++ {
		go func() {
			c.Add(1)
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println(c.Count())

	c.Reset()
	fmt.Println(c.Count())

}
