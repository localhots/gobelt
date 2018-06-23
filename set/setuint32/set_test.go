package impl

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/google/go-cmp/cmp"
)

const (
	One   uint32 = 1
	Two   uint32 = 2
	Three uint32 = 3
	Four  uint32 = 4
	Five  uint32 = 5
)

func TestNew(t *testing.T) {
	s := New(One, Two, Three)
	if s == nil {
		t.Fatal("Set is nil")
	}
	if s.Len() != 3 {
		t.Errorf("Expected set to contain 3 items, got %d", s.Len())
	}
	for _, item := range []uint32{One, Two, Three} {
		if ok := s.Has(item); !ok {
			t.Errorf("Set is expected to contain item %q", item)
		}
	}
}

func TestAdd(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Errorf("Expected set to be empty, got %d items", s.Len())
	}
	s.Add(One)
	s.Add(Two, Three)
	for _, item := range []uint32{One, Two, Three} {
		if ok := s.Has(item); !ok {
			t.Errorf("Set is expected to contain item %q", item)
		}
	}
}

func TestRemove(t *testing.T) {
	s := New(One, Two, Three, Four, Five)
	s.Remove(One, Two)
	s.Remove(Three)
	if s.Len() != 2 {
		t.Errorf("Expected set to contain 2 items, got %d", s.Len())
	}
	for _, item := range []uint32{One, Two, Three} {
		if ok := s.Has(item); ok {
			t.Errorf("Set is expected to not contain item %q", item)
		}
	}
	for _, item := range []uint32{Four, Five} {
		if ok := s.Has(item); !ok {
			t.Errorf("Set is expected to contain item %q", item)
		}
	}
}

func TestHas(t *testing.T) {
	s := New(One, Two)
	table := map[uint32]bool{
		One:   true,
		Two:   true,
		Three: false,
		Four:  false,
		Five:  false,
	}
	for v, exp := range table {
		if res := s.Has(v); res != exp {
			t.Errorf("Item: %v, In: %v, Expected: %v", v, res, exp)
		}
	}
}

func TestLen(t *testing.T) {
	table := map[*Set]int{
		New():                                  0,
		New(One):                               1,
		New(Two, Three):                        2,
		New(One, Two, Three, Four, Five, Five): 5,
	}
	for s, exp := range table {
		if res := s.Len(); res != exp {
			t.Errorf("Expected set %s to have length %d, got %d", s, exp, res)
		}
	}
}

func TestSlice(t *testing.T) {
	s := New(One, Two, Three)
	out := s.Slice()
	exp := []uint32{One, Two, Three}
	ignoreOrder := cmpopts.SortSlices(func(a, b uint32) bool {
		return a < b
	})
	if !cmp.Equal(exp, out, ignoreOrder) {
		t.Errorf("Retured slice does not match: %s", cmp.Diff(exp, out))
	}
}

func TestSortedSlice(t *testing.T) {
	s := New(One, Two, Three)
	out := s.SortedSlice()
	exp := []uint32{One, Two, Three}
	if !cmp.Equal(exp, out) {
		t.Errorf("Retured slice does not match: %s", cmp.Diff(exp, out))
	}
}
