package util

import (
	"fmt"
	"testing"

	"github.com/unk1ndled/draw/util"
)

// TestMap tests the Map function.
func TestMap(t *testing.T) {
	tests := []struct {
		x, in_min, in_max, out_min, out_max, expected float64
	}{
		// Positive values
		{5, 0, 10, 0, 100, 50},
		{2.5, 0, 10, 0, 100, 25},
		{0, 0, 10, 0, 100, 0},
		{10, 0, 10, 0, 100, 100},
		{5, 0, 10, 100, 200, 150},
		{5, 1, 11, 0, 100, 40},
		{7.5, 0, 10, 0, 1, 0.75},

		// Negative values
		{-5, -10, 0, 0, 100, 50},
		{-2.5, -10, 0, 0, 100, 75},
		{-10, -10, 0, 0, 100, 0},
		{-5, -10, 0, 100, 200, 150},
		{-5, -10, 0, -100, 0, -50},

		// Mixed positive and negative ranges
		{5, -10, 10, 0, 100, 75},
		{0, -10, 10, 0, 100, 50},
		{-5, -10, 10, 0, 100, 25},
		{-10, -10, 10, 0, 100, 0},
		{10, -10, 10, 0, 100, 100},
		{0, -10, 10, -50, 50, 0},
		{-10, -10, 10, -50, 50, -50},
		{10, -10, 10, -50, 50, 50},

		// Other random tests
		{3, 0, 6, 0, 1, 0.5},
		{7, 0, 14, 0, 100, 50},
		{15, 10, 20, 100, 200, 150},
		{1, 0, 1, 0, 100, 100},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Map(%v, %v, %v, %v, %v)", tt.x, tt.in_min, tt.in_max, tt.out_min, tt.out_max), func(t *testing.T) {
			result := util.Map(tt.x, tt.in_min, tt.in_max, tt.out_min, tt.out_max)
			if result != tt.expected {
				t.Errorf("Map(%v, %v, %v, %v, %v) = %v; want %v", tt.x, tt.in_min, tt.in_max, tt.out_min, tt.out_max, result, tt.expected)
			}
		})
	}
}
