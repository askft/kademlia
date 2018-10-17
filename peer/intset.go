package peer

import (
	"sort"
	"sync"
)

// OptionalTODO (but totally unnecessary): union, difference, intersection

// IntSet is a set of `int` values.
type IntSet struct {
	set    map[int]int
	values []int
	mutex  sync.Mutex
}

// NewIntSet returns a new IntSet.
func NewIntSet() *IntSet {
	return &IntSet{
		set:    make(map[int]int),
		values: make([]int, 0),
		mutex:  sync.Mutex{},
	}
}

// Add adds `v` to `s`.
func (s *IntSet) Add(v int) {
	if _, ok := s.set[v]; ok {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.set[v] = len(s.set)
	s.values = append(s.values, v)
}

// AddMany adds all values in `vs` to `s`.
func (s *IntSet) AddMany(vs []int) {
	for _, v := range vs {
		s.Add(v)
	}
}

// Remove removes `v` from `s`.
func (s *IntSet) Remove(v int) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	i := s.set[v]
	delete(s.set, v)
	s.values = append(s.values[:i], s.values[i+1:]...)
}

// Has returns `true` if `v` is in `s`, otherwise false.
func (s *IntSet) Has(v int) bool {
	_, ok := s.set[v]
	return ok
}

// Values returns the a copy of the set values
// sorted by insertion time (oldest first).
func (s *IntSet) Values() []int {
	return append([]int(nil), s.values...)
}

// Sorted returns the values in the set sorted by value, ascending.
func (s *IntSet) Sorted() []int {
	values := s.Values()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sort.Ints(values)
	return values
}

// SortedReverse returns the values in the set sorted by value, descending.
func (s *IntSet) SortedReverse() []int {
	values := s.Values()
	s.mutex.Lock()
	defer s.mutex.Unlock()
	sort.Sort(sort.Reverse(sort.IntSlice(values)))
	return values
}
