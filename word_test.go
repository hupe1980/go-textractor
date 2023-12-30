package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestWord(t *testing.T) {
	// Create a sample line for testing
	line := &Line{
		base: base{
			id:          "line-1",
			confidence:  0.9,
			blockType:   types.BlockTypeLine,
			boundingBox: &BoundingBox{},
			page:        &Page{},
			raw:         types.Block{},
		},
		words: []*Word{
			{
				base: base{
					id:          "word-1",
					confidence:  0.9,
					blockType:   types.BlockTypeWord,
					boundingBox: &BoundingBox{},
					page:        &Page{},
					raw:         types.Block{},
				},
				text:      "Hello",
				textType:  types.TextTypePrinted,
				line:      nil, // Link to the line is set later for testing
				tableCell: nil, // Link to the table cell is set later for testing
			},
			{
				base: base{
					id:          "word-2",
					confidence:  0.9,
					blockType:   types.BlockTypeWord,
					boundingBox: &BoundingBox{},
					page:        &Page{},
					raw:         types.Block{},
				},
				text:      "World",
				textType:  types.TextTypeHandwriting,
				line:      nil, // Link to the line is set later for testing
				tableCell: nil, // Link to the table cell is set later for testing
			},
		},
	}

	// Set the line reference in each word
	for _, word := range line.words {
		word.line = line
	}

	// Test Word methods
	t.Run("Word_Text", func(t *testing.T) {
		assert.Equal(t, "Hello", line.words[0].Text())
		assert.Equal(t, "World", line.words[1].Text())
	})

	t.Run("Word_TextType", func(t *testing.T) {
		assert.Equal(t, types.TextTypePrinted, line.words[0].TextType())
		assert.Equal(t, types.TextTypeHandwriting, line.words[1].TextType())
	})

	t.Run("Word_IsPrinted", func(t *testing.T) {
		assert.True(t, line.words[0].IsPrinted())
		assert.True(t, line.words[1].IsHandwriting())
	})

	t.Run("Word_IsHandwriting", func(t *testing.T) {
		assert.False(t, line.words[0].IsHandwriting())
		assert.False(t, line.words[1].IsPrinted())
	})
}
