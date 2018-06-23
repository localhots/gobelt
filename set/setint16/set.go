/*******************************************************************************
THIS FILE WAS AUTOMATICALLY GENERATED. DO NOT EDIT!
*******************************************************************************/
package impl

import (
	"fmt"
	"sort"
)

// Set is a set of int16.
type Set struct {
	items map[int16]struct{}
}

// New creates a new int16 set.
func New(items ...int16) *Set {
	s := &Set{items: make(map[int16]struct{}, len(items))}
	s.Add(items...)
	return s
}

// Add adds given items to the set.
func (s *Set) Add(items ...int16) *Set {
	for _, item := range items {
		s.items[item] = struct{}{}
	}
	return s
}

// Remove delete given items from the set.
func (s *Set) Remove(items ...int16) *Set {
	for _, item := range items {
		delete(s.items, item)
	}
	return s
}

// Has returns true if all the given items are included in the set.
func (s *Set) Has(items ...int16) bool {
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
func (s *Set) Slice() []int16 {
	sl := make([]int16, len(s.items))
	i := 0
	for item := range s.items {
		sl[i] = item
		i++
	}
	return sl
}

// SortedSlice returns items of the set as a slice sorted ascending.
func (s *Set) SortedSlice() []int16 {
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
