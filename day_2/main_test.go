package main

import "testing"

func equal(a []int, b []int) bool {
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestEliminateIndex(t *testing.T) {
	cases := []struct {
		Name      string
		Line      []int
		Eliminate int
		Expected  []int
	}{
		{
			Name: "0th Item", Line: []int{1, 2, 3, 4, 5},
			Eliminate: 0,
			Expected:  []int{2, 3, 4, 5},
		},
		{
			Name:      "1st Item",
			Line:      []int{1, 2, 3, 4, 5},
			Eliminate: 1,
			Expected:  []int{1, 3, 4, 5},
		},
		{
			Name:      "Last Item",
			Line:      []int{1, 2, 3, 4, 5},
			Eliminate: 4,
			Expected:  []int{1, 2, 3, 4},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			eliminated := eliminateIndex(c.Line, c.Eliminate)

			if len(eliminated) != len(c.Expected) {
				t.Errorf("Expected %v, got %v", c.Expected, eliminated)
			}

			if !equal(eliminated, c.Expected) {
				t.Errorf("Expected %v, got %v", c.Expected, eliminated)
			}
		})
	}
}
