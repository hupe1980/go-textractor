package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	t.Run("Integegers", func(t *testing.T) {
		result := Sum([]int{1, 2, 3})
		assert.Equal(t, 6, result)
	})

	t.Run("Floats", func(t *testing.T) {
		result := Sum([]float32{1.5, 2.5, 3.5})
		assert.Equal(t, float32(7.5), result)
	})

	t.Run("EmptySlice", func(t *testing.T) {
		result := Sum([]int{})
		assert.Equal(t, 0, result)
	})
}

func TestMean(t *testing.T) {
	t.Run("Integegers", func(t *testing.T) {
		result := Mean([]int{1, 2, 3, 4, 5})
		assert.Equal(t, float32(3), result)
	})

	t.Run("Floats", func(t *testing.T) {
		result := Mean([]float32{1.5, 2.5, 3.5})
		assert.Equal(t, float32(2.5), result)
	})

	t.Run("EmptySlice", func(t *testing.T) {
		result := Mean([]int{})
		assert.Equal(t, float32(0), result)
	})
}

func TestValues(t *testing.T) {
	tests := []struct {
		name   string
		input  map[string]int
		output []int
	}{
		{"IntMap", map[string]int{"a": 1, "b": 2, "c": 3}, []int{1, 2, 3}},
		{"EmptyMap", map[string]int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Values(tt.input)
			assert.ElementsMatch(t, tt.output, result)
		})
	}
}
