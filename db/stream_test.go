package db

import (
	"testing"
	"unsafe"
)

func TestStream(t *testing.T) {
	s1 := []byte("Hello")

	s2, err := format(*(*string)(unsafe.Pointer(&s1)), "string")
	if err != nil {
		t.Logf("%v", err)
		return
	}
	t.Logf("s1=%s, %p", s1, s1)
	t.Logf("s2=%s, %p", s2, &s2)
}
