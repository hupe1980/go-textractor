package textractor

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// SelectionElement represents a selectable element in the document.
type SelectionElement struct {
	content
}

// NewSelectionElement creates a new SelectionElement instance.
func NewSelectionElement(block types.Block) *SelectionElement {
	return &SelectionElement{
		content: content{block},
	}
}

// Status returns the selection status of the element.
func (se *SelectionElement) Status() types.SelectionStatus {
	return se.Block().SelectionStatus
}

// IsSelected checks if the element is selected.
func (se *SelectionElement) IsSelected() bool {
	return se.Status() == types.SelectionStatusSelected
}

// Word represents a word in the document.
type Word struct {
	content
}

// NewWord creates a new Word instance.
func NewWord(block types.Block) *Word {
	word := &Word{
		content: content{block},
	}

	return word
}

// TextType returns the text type of the word.
func (w *Word) TextType() types.TextType {
	return w.Block().TextType
}

// IsPrinted checks if the word is printed text.
func (w *Word) IsPrinted() bool {
	return w.TextType() == types.TextTypePrinted
}

// IsHandwriting checks if the word is handwriting.
func (w *Word) IsHandwriting() bool {
	return w.TextType() == types.TextTypeHandwriting
}

// Text returns the text content of the word.
func (w *Word) Text() string {
	return aws.ToString(w.Block().Text)
}

// Line represents a line of text in the document.
type Line struct {
	words []*Word
	content
}

// NewLine creates a new Line instance.
func NewLine(block types.Block, blockMap map[string]types.Block) *Line {
	line := &Line{
		content: content{block},
	}

	for _, r := range block.Relationships {
		if r.Type == types.RelationshipTypeChild {
			for _, i := range r.Ids {
				w := blockMap[i]
				if w.BlockType == types.BlockTypeWord {
					line.words = append(line.words, NewWord(w))
				}
			}
		}
	}

	return line
}

// Words returns the words in the line.
func (l *Line) Words() []*Word {
	return l.words
}

// Text returns the text content of the line.
func (l *Line) Text() string {
	return aws.ToString(l.Block().Text)
}

// Content is an interface for document content elements.
type Content interface {
	ID() string
	Confidence() float32
	Geometry() *Geometry
	Block() types.Block
}

// content is a common struct embedding a types.Block to provide basic functionality.
type content struct {
	block types.Block
}

// ID returns the ID of the content.
func (c *content) ID() string {
	return aws.ToString(c.block.Id)
}

// Confidence returns the confidence level of the content.
func (c *content) Confidence() float32 {
	return aws.ToFloat32(c.block.Confidence)
}

// Geometry returns the geometry information of the content.
func (c *content) Geometry() *Geometry {
	return NewGeometry(c.block.Geometry)
}

// Block returns the underlying types.Block of the content.
func (c *content) Block() types.Block {
	return c.block
}
