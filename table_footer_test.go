package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestTableFooter(t *testing.T) {
	// Mock data
	word1 := &Word{text: "First", textType: types.TextTypePrinted}
	word2 := &Word{text: "Footer", textType: types.TextTypePrinted}
	words := []*Word{word1, word2}

	// Create a TableFooter instance
	tableFooter := &TableFooter{
		base: base{
			id:         "footer-1",
			confidence: 0.95,
			blockType:  types.BlockTypeTableFooter,
		},
		words: words,
	}

	// Test Words method
	assert.Equal(t, words, tableFooter.Words(), "Words method should return the correct words")

	// Test Text method
	expectedText := "First Footer"
	assert.Equal(t, expectedText, tableFooter.Text(), "Text method should return the concatenated text")
}
