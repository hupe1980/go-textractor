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
		assert.Equal(t, float64(3), result)
	})

	t.Run("Floats32", func(t *testing.T) {
		result := Mean([]float32{1.5, 2.5, 3.5})
		assert.Equal(t, float64(2.5), result)
	})

	t.Run("Floats64", func(t *testing.T) {
		result := Mean([]float64{1.5, 2.5, 3.5})
		assert.Equal(t, float64(2.5), result)
	})

	t.Run("EmptySlice", func(t *testing.T) {
		result := Mean([]int{})
		assert.Equal(t, float64(0), result)
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

func TestConcatenate(t *testing.T) {
	t.Run("ConcatenateInts", func(t *testing.T) {
		ints1 := []int{1, 2, 3}
		ints2 := []int{4, 5, 6}
		expectedInts := []int{1, 2, 3, 4, 5, 6}
		resultInts := Concatenate(ints1, ints2)
		assert.Equal(t, expectedInts, resultInts, "Concatenate of ints failed")
	})

	t.Run("ConcatenateStrings", func(t *testing.T) {
		strings1 := []string{"a", "b", "c"}
		strings2 := []string{"d", "e", "f"}
		expectedStrings := []string{"a", "b", "c", "d", "e", "f"}
		resultStrings := Concatenate(strings1, strings2)
		assert.Equal(t, expectedStrings, resultStrings, "Concatenate of strings failed")
	})

	t.Run("ConcatenateEmptySlices", func(t *testing.T) {
		emptyInts := []int{}
		emptyStrings := []string{}
		emptyMixed := []interface{}{}
		resultEmptyInts := Concatenate(emptyInts)
		resultEmptyStrings := Concatenate(emptyStrings)
		resultEmptyMixed := Concatenate(emptyMixed)
		assert.Empty(t, resultEmptyInts, "Concatenate of empty ints failed")
		assert.Empty(t, resultEmptyStrings, "Concatenate of empty strings failed")
		assert.Empty(t, resultEmptyMixed, "Concatenate of empty mixed types failed")
	})
}
