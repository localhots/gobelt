/*******************************************************************************
THIS FILE WAS AUTOMATICALLY GENERATED. DO NOT EDIT!
*******************************************************************************/
package impl

import (
	"fmt"
	"sort"
)

// Set is a set of uint64.
type Set struct {
	items map[uint64]struct{}
}

// New creates a new uint64 set.
func New(items ...uint64) *Set {
	s := &Set{items: make(map[uint64]struct{}, len(items))}
	s.Add(items...)
	return s
}

// Add adds given items to the set.
func (s *Set) Add(items ...uint64) *Set {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
	return s
}

// Remove delete given items from the set.
func (s *Set) Remove(items ...uint64) *Set {
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

// Has returns true if all the given items are included in the set.
func (s *Set) Has(items ...uint64) bool {
	for _, item := range items {
		if _, ok := s.items[item]; !ok {
			return false
		}
	}
	return true
}

// Len returns the size of the set.
func (s *Set) Len() int {
	return len(s.items)
}

// Slice returns items of the set as a slice.
func (s *Set) Slice() []uint64 {
	sl := make([]uint64, len(s.items))
	i := 0
	for item := range s.items {
		sl[i] = item
		i++
	}
	return sl
}

// SortedSlice returns items of the set as a slice sorted ascending.
func (s *Set) SortedSlice() []uint64 {
	ss := s.Slice()
	sort.Slice(ss, func(i, j int) bool {
		return ss[i] < ss[j]
	})
	return ss
}

// String implements fmt.Stringer interface.
func (s *Set) String() string {
	return fmt.Sprintf("[%v]", s.Slice())
}
