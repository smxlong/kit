package mapindex

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// pair is a simple key-value pair.
type pair struct {
	key   string
	value int
}

// Test_MapIndex_Basic_Usage tests that we can create a map index and get the
// keys and values from it, and that the index has the correct values at the
// correct keys.
func Test_MapIndex_Basic_Usage(t *testing.T) {
	slice := []pair{{"one", 1}, {"two", 2}, {"three", 3}, {"four", 4}, {"five", 5}, {"six", 6}, {"seven", 7}, {"eight", 8}, {"nine", 9}, {"ten", 10}}
	m := New(slice, func(p pair) string { return p.key })
	keys := Keys(m)
	values := Values(m)
	compareAnyOrder(t, Map(slice, func(p pair) string { return p.key }), keys)
	compareAnyOrder(t, Map(slice, func(p pair) int { return p.value }), Map(values, func(p pair) int { return p.value }))
	for _, p := range slice {
		if m[p.key] != p {
			t.Errorf("expected %v, got %v", p, m[p.key])
		}
	}
}

// compareAnyOrder compares two slices, ignoring the order of the elements.
func compareAnyOrder[T comparable](t *testing.T, expected, actual []T) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
	for _, e := range expected {
		found := false
		for _, a := range actual {
			if e == a {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	}
}

// Test_Filter tests the Filter function.
func Test_Filter(t *testing.T) {
	// Test cases with three input values and one expected result.
	tc := []struct {
		in  []int
		out []int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, []int{2, 4, 6, 8, 10}},
		{[]int{1, 3, 5, 7, 9}, nil},
		{[]int{2, 4, 6, 8, 10}, []int{2, 4, 6, 8, 10}},
	}

	for _, c := range tc {
		require.Equal(t, c.out, Filter(c.in, func(x int) bool { return x%2 == 0 }))
	}
}
