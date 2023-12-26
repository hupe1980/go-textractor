package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	t.Run("Parse doc with liens and words", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(&AnalyzeDocumentPage{Blocks: td.Blocks})
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, 22, len(doc.Pages()[0].Lines()))

		wc := 0
		for _, l := range doc.Pages()[0].Lines() {
			wc += len(l.Words())
		}
		assert.Equal(t, 53, wc)
	})
}

func TestSelectionElement(t *testing.T) {
	block := types.Block{
		Id:              aws.String("123"),
		Confidence:      aws.Float32(99.123),
		SelectionStatus: types.SelectionStatusSelected,
		Geometry: &types.Geometry{
			BoundingBox: &types.BoundingBox{},
			Polygon:     []types.Point{},
		},
	}

	selectionElement := NewSelectionElement(block)

	assert.Equal(t, "123", selectionElement.ID())
	assert.Equal(t, float32(99.123), selectionElement.Confidence())
	assert.NotNil(t, selectionElement.Geometry())
	assert.Equal(t, block, selectionElement.Block())
	assert.Equal(t, types.SelectionStatusSelected, selectionElement.Status())
	assert.True(t, selectionElement.IsSelected())
}

func TestWord(t *testing.T) {
	block := types.Block{
		Id:         aws.String("456"),
		TextType:   types.TextTypePrinted,
		Confidence: aws.Float32(99.123),
		Text:       aws.String("example"),
		Geometry: &types.Geometry{
			BoundingBox: &types.BoundingBox{},
			Polygon:     []types.Point{},
		},
	}

	word := NewWord(block)

	assert.Equal(t, "456", word.ID())
	assert.Equal(t, float32(99.123), word.Confidence())
	assert.NotNil(t, word.Geometry())
	assert.Equal(t, block, word.Block())
	assert.Equal(t, types.TextTypePrinted, word.TextType())
	assert.True(t, word.IsPrinted())
	assert.False(t, word.IsHandwriting())
	assert.Equal(t, "example", word.Text())
}

func TestLine(t *testing.T) {
	block := types.Block{
		Id:         aws.String("789"),
		Confidence: aws.Float32(99.123),
		Relationships: []types.Relationship{
			{
				Type: types.RelationshipTypeChild,
				Ids:  []string{"101", "102", "103", "104"},
			},
		},
		Geometry: &types.Geometry{
			BoundingBox: &types.BoundingBox{},
			Polygon:     []types.Point{},
		},
		Text: aws.String("This is a line."),
	}

	blockMap := map[string]types.Block{
		"101": {Id: aws.String("101"), BlockType: types.BlockTypeWord, Text: aws.String("This")},
		"102": {Id: aws.String("102"), BlockType: types.BlockTypeWord, Text: aws.String("is")},
		"103": {Id: aws.String("103"), BlockType: types.BlockTypeWord, Text: aws.String("a")},
		"104": {Id: aws.String("104"), BlockType: types.BlockTypeWord, Text: aws.String("line.")},
	}

	line := NewLine(block, blockMap)

	assert.Equal(t, "789", line.ID())
	assert.Equal(t, float32(99.123), line.Confidence())
	assert.NotNil(t, line.Geometry())
	assert.Equal(t, block, line.Block())
	assert.Equal(t, "This is a line.", line.Text())
	assert.Len(t, line.Words(), 4)

	// Testing the words within the line.
	word1 := line.Words()[0]
	word2 := line.Words()[1]
	word3 := line.Words()[2]
	word4 := line.Words()[3]

	assert.Equal(t, "101", word1.ID())
	assert.Equal(t, "This", word1.Text())

	assert.Equal(t, "102", word2.ID())
	assert.Equal(t, "is", word2.Text())

	assert.Equal(t, "103", word3.ID())
	assert.Equal(t, "a", word3.Text())

	assert.Equal(t, "104", word4.ID())
	assert.Equal(t, "line.", word4.Text())
}
