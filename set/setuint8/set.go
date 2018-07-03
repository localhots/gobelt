/*******************************************************************************
THIS FILE WAS AUTOMATICALLY GENERATED. DO NOT EDIT!
*******************************************************************************/
package setuint8

import (
	"fmt"
	"sort"
)

// Set is a set of uint8.
type Set struct {
	items map[uint8]struct{}
}

// New creates a new uint8 set.
func New(items ...uint8) *Set {
	s := &Set{items: make(map[uint8]struct{}, len(items))}
	s.Add(items...)
	return s
}

// Add adds given items to the set.
func (s *Set) Add(items ...uint8) *Set {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
	return s
}

// Remove delete given items from the set.
func (s *Set) Remove(items ...uint8) *Set {
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

// Has returns true if all the given items are included in the set.
func (s *Set) Has(items ...uint8) bool {
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
func (s *Set) Slice() []uint8 {
	sl := make([]uint8, len(s.items))
	i := 0
	for item := range s.items {
		sl[i] = item
		i++
	}
	return sl
}

// SortedSlice returns items of the set as a slice sorted ascending.
func (s *Set) SortedSlice() []uint8 {
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
