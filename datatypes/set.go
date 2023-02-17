package datatypes

import "sync"

// Set set
type Set struct {
	m   map[string]struct{}
	mux sync.RWMutex
}

// NewSet new set
func NewSet(vs ...string) *Set {
	s := Set{
		m: make(map[string]struct{}),
	}
	s.Append(vs...)
	return &s
}

// Append append
func (s *Set) Append(vs ...string) *Set {
	s.mux.Lock()
	defer s.mux.Unlock()
	for _, v := range vs {
		s.m[v] = struct{}{}
	}
	return s
}

// Len len
func (s *Set) Len() int {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return len(s.m)
}

// List list
func (s *Set) List() []string {
	sList := make([]string, 0, s.Len())
	s.mux.RLock()
	defer s.mux.RUnlock()
	for v := range s.m {
		sList = append(sList, v)
	}
	return sList
}

// Contains contains
func (s *Set) Contains(v string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, ok := s.m[v]
	return ok
}
