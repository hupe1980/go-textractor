package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTableTitle(t *testing.T) {
	t.Run("Words", func(t *testing.T) {
		// Create a TableTitle with some words
		words := []*Word{
			{base: base{id: "1", confidence: 0.9, blockType: "Word"}, text: "Hello"},
			{base: base{id: "2", confidence: 0.8, blockType: "Word"}, text: "World"},
		}
		tableTitle := &TableTitle{base: base{id: "tableTitle", confidence: 0.95, blockType: "TableTitle"}, words: words}

		// Test the Words method
		result := tableTitle.Words()
		assert.Equal(t, words, result)
	})

	t.Run("Text", func(t *testing.T) {
		// Create a TableTitle with some words
		words := []*Word{
			{base: base{id: "1", confidence: 0.9, blockType: "Word"}, text: "Hello"},
			{base: base{id: "2", confidence: 0.8, blockType: "Word"}, text: "World"},
		}
		tableTitle := &TableTitle{base: base{id: "tableTitle", confidence: 0.95, blockType: "TableTitle"}, words: words}

		// Test the Text method
		result := tableTitle.Text()
		expectedText := "Hello World"
		assert.Equal(t, expectedText, result)
	})
}
