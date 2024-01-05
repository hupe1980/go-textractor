package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeLevenshteinDistance(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		str1     string
		str2     string
		expected int
	}{
		{"kitten", "sitting", 3},      // Example from Wikipedia
		{"", "hello", 5},              // Empty string to a non-empty string
		{"abc", "abc", 0},             // Same strings
		{"abc", "def", 3},             // Completely different strings
		{"abcdef", "abefcd", 4},       // Random strings
		{"hello", "hola", 3},          // Different lengths
		{"intention", "execution", 5}, // Another example from Wikipedia
		{"kangaroo", "koal", 6},       // Different lengths
		{"k", "k", 0},                 // Single-character strings
		{"", "", 0},                   // Empty strings
		{"abcd", "", 4},               // Empty string to a non-empty string
	}

	for _, testCase := range testCases {
		result := ComputeLevenshteinDistance(testCase.str1, testCase.str2)
		assert.Equal(testCase.expected, result, "Error in test case: %s, %s", testCase.str1, testCase.str2)
	}
}

func BenchmarkComputeLevenshteinDistance(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// Adjust the input strings as needed for your benchmark
		str1 := "kitten"
		str2 := "sitting"
		ComputeLevenshteinDistance(str1, str2)
	}
}
