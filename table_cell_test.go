package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestTableMergedCell(t *testing.T) {
	t.Run("Words", func(t *testing.T) {
		// Create sample words
		word1 := &Word{text: "Hello", textType: types.TextTypePrinted}
		word2 := &Word{text: "World", textType: types.TextTypePrinted}

		// Create a TableMergedCell with sample cells
		cell1 := &TableCell{words: []*Word{word1}}
		cell2 := &TableCell{words: []*Word{word2}}
		mergedCell := &TableMergedCell{cells: []*TableCell{cell1, cell2}}

		// Call the Words method and check if it returns the expected words
		result := mergedCell.Words()
		assert.Equal(t, []*Word{word1, word2}, result, "Words mismatch")
	})

	t.Run("Text", func(t *testing.T) {
		// Create sample words
		word1 := &Word{text: "Hello", textType: types.TextTypePrinted}
		word2 := &Word{text: "World", textType: types.TextTypePrinted}

		// Create a TableMergedCell with sample cells
		cell1 := &TableCell{words: []*Word{word1}}
		cell2 := &TableCell{words: []*Word{word2}}
		mergedCell := &TableMergedCell{cells: []*TableCell{cell1, cell2}}

		// Call the Text method and check if it returns the expected text
		result := mergedCell.Text()
		assert.Equal(t, "Hello World", result, "Text mismatch")
	})
}

func TestTableCell(t *testing.T) {
	t.Run("Words", func(t *testing.T) {
		// Create sample words
		word1 := &Word{text: "Hello", textType: types.TextTypePrinted}
		word2 := &Word{text: "World", textType: types.TextTypePrinted}

		// Create a TableCell with sample words
		cell := &TableCell{words: []*Word{word1, word2}}

		// Call the Words method and check if it returns the expected words
		result := cell.Words()
		assert.Equal(t, []*Word{word1, word2}, result, "Words mismatch")
	})

	t.Run("SelectionElement", func(t *testing.T) {
		// Create a sample SelectionElement
		selectionElement := &SelectionElement{status: types.SelectionStatusSelected}

		// Create a TableCell with the sample SelectionElement
		cell := &TableCell{selectionElement: selectionElement}

		// Call the SelectionElement method and check if it returns the expected SelectionElement
		result := cell.SelectionElement()
		assert.Equal(t, selectionElement, result, "SelectionElement mismatch")
	})

	t.Run("Text", func(t *testing.T) {
		// Create sample words
		word1 := &Word{text: "Hello", textType: types.TextTypePrinted}
		word2 := &Word{text: "World", textType: types.TextTypePrinted}

		// Create a TableCell with sample words
		cell := &TableCell{words: []*Word{word1, word2}}

		// Call the Text method and check if it returns the expected text
		result := cell.Text()
		assert.Equal(t, "Hello World", result, "Text mismatch")
	})
}
