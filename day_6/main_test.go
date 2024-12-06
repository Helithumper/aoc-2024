package main

import (
	"reflect"
	"testing"
)

func TestFilterValidOrders(t *testing.T) {
	tests := []struct {
		name         string
		lines        []string
		prePostRules map[int][]int
		postPreRules map[int][]int
		expected     [][]int
	}{
		{
			name:         "Basic test",
			lines:        []string{"1,2,3"},
			prePostRules: map[int][]int{1: {2}, 2: {3}},
			postPreRules: map[int][]int{2: {1}, 3: {2}},
			expected:     [][]int{{1, 2, 3}},
		}, {
			name:         "Basic test 2",
			lines:        []string{"1,2,3", "1,3,2"},
			prePostRules: map[int][]int{1: {2, 3}, 2: {3}},
			postPreRules: map[int][]int{2: {1}, 3: {1, 2}},
			expected:     [][]int{{1, 2, 3}},
		}, {
			name:         "Invalid PrePostRules",
			lines:        []string{"1,2,3", "1,3,2", "2,3,1"},
			prePostRules: map[int][]int{1: {2, 3}, 2: {3}, 3: {1}},
			postPreRules: map[int][]int{1: {3}, 2: {1}, 3: {1, 2}},
			expected:     [][]int{},
		}, {
			name:         "Invalid PostPreRules",
			lines:        []string{"1,2,3", "1,3,2", "2,3,1"},
			prePostRules: map[int][]int{1: {2, 3}, 2: {3}},
			postPreRules: map[int][]int{1: {3}, 2: {1}, 3: {1, 2}},
			expected:     [][]int{},
		}, {
			name:         "Irrelevant rules",
			lines:        []string{"1,2,3", "1,3,2", "2,3,1"},
			prePostRules: map[int][]int{7: {2, 3}, 12: {3}, 13: {1}},
			postPreRules: map[int][]int{9: {3}, 12: {1}, 13: {1, 2}},
			expected:     [][]int{{1, 2, 3}, {1, 3, 2}, {2, 3, 1}},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := filterValidOrders(tc.lines, tc.prePostRules, tc.postPreRules)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Errorf("Expected %v, got %v", tc.expected, got)
			}
		})
	}
}
