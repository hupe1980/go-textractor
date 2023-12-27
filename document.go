package textractor

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor/internal"
)

// Document represents a Textract document containing pages.
type Document struct {
	blockMap map[string]types.Block
	pages    []*Page
}

// NewDocument creates a new Document instance using response pages from Textract.
func NewDocument(documentBlocks ...[]types.Block) *Document {
	doc := &Document{
		blockMap: make(map[string]types.Block),
	}

	var (
		currentPageBlock   *types.Block
		currentPageContent []types.Block
	)

	for _, blocks := range documentBlocks {
		for i, b := range blocks {
			doc.blockMap[aws.ToString(b.Id)] = b

			if b.BlockType == types.BlockTypePage {
				if currentPageBlock != nil {
					doc.pages = append(doc.pages, NewPage(*currentPageBlock, currentPageContent, doc.blockMap))
				}

				currentPageBlock = &blocks[i]
				currentPageContent = make([]types.Block, 0)
				currentPageContent = append(currentPageContent, b)
			} else {
				currentPageContent = append(currentPageContent, b)
			}
		}

		if currentPageBlock != nil {
			doc.pages = append(doc.pages, NewPage(*currentPageBlock, currentPageContent, doc.blockMap))
		}
	}

	return doc
}

// BlockByID retrieves a block by its ID.
func (doc *Document) BlockByID(id string) *types.Block {
	b, ok := doc.blockMap[id]
	if !ok {
		return nil
	}

	return &b
}

// PageNumber retrieves a page by its page number.
func (doc *Document) PageNumber(n int) *Page {
	if n < 1 || n > len(doc.pages) {
		panic(fmt.Sprintf("number %d must be between 1 and %d", n, len(doc.pages)))
	}

	return doc.pages[n-1]
}

// Pages returns all pages in the document.
func (doc *Document) Pages() []*Page {
	return doc.pages
}

// PageCount returns the total number of pages in the document.
func (doc *Document) PageCount() int {
	return len(doc.pages)
}

// Page represents a page in the document.
type Page struct {
	block      types.Block
	blocks     []types.Block
	lines      []*Line
	form       *Form
	tables     []*Table
	signatures []*Signature
	queries    Queries
}

// NewPage creates a new Page instance using Textract page blocks and a block map.
func NewPage(pageBlock types.Block, blocks []types.Block, blockMap map[string]types.Block) *Page {
	page := &Page{
		block:  pageBlock,
		blocks: blocks,
		form:   NewForm(),
	}

	for _, b := range blocks {
		switch b.BlockType {
		case types.BlockTypeLine:
			page.lines = append(page.lines, NewLine(b, blockMap))
		case types.BlockTypeTable:
			page.tables = append(page.tables, NewTable(b, blockMap))
		case types.BlockTypeKeyValueSet:
			if slices.Contains(b.EntityTypes, types.EntityTypeKey) {
				f := NewField(b, blockMap)
				if f.Key() != nil {
					page.form.AddField(f)
				}
			}
		case types.BlockTypeQuery:
			page.queries = append(page.queries, NewQuery(b, blockMap))
		case types.BlockTypeSignature:
			page.signatures = append(page.signatures, NewSignature(b))
		default: // TODO logging?
		}
	}

	return page
}

// ID returns the ID of the page block.
func (p *Page) ID() string {
	return aws.ToString(p.block.Id)
}

// Blocks returns all blocks in the page.
func (p *Page) Blocks() []types.Block {
	return p.blocks
}

// Geometry returns the geometry of the page.
func (p *Page) Geometry() *Geometry {
	return NewGeometry(p.block.Geometry)
}

// Text returns the concatenated text from all lines in the page.
func (p *Page) Text() string {
	texts := make([]string, len(p.lines))
	for i, l := range p.lines {
		texts[i] = l.Text()
	}

	return strings.Join(texts, "\n")
}

// Form returns the form information on the page.
func (p *Page) Form() *Form {
	return p.form
}

// Signatures returns the signatures on the page.
func (p *Page) Signatures() []*Signature {
	return p.signatures
}

// Queries returns the queries for the page.
func (p *Page) Queries() Queries {
	return p.queries
}

// TableCount returns the total number of tables in the page.
func (p *Page) TableCount() int {
	return len(p.tables)
}

// Tables returns all tables in the page.
func (p *Page) Tables() []*Table {
	return p.tables
}

// TableAtIndex returns the table at the specified index.
func (p *Page) TableAtIndex(i int) *Table {
	if i < 0 || i >= len(p.tables) {
		panic(fmt.Sprintf("index %d must be > 0 and < %d", i, len(p.tables)))
	}

	return p.tables[i]
}

// LineCount returns the total number of lines in the page.
func (p *Page) LineCount() int {
	return len(p.lines)
}

// Lines returns all lines in the page.
func (p *Page) Lines() []*Line {
	return p.lines
}

// LineAtIndex returns the line at the specified index.
func (p *Page) LineAtIndex(i int) *Line {
	if i < 0 || i >= len(p.lines) {
		panic(fmt.Sprintf("index %d must be > 0 and < %d", i, len(p.lines)))
	}

	return p.lines[i]
}

// OCRConfidence represents the OCR confidence scores, including mean, max, and min values.
type OCRConfidence struct {
	mean float32
	max  float32
	min  float32
}

// NewOCRConfidenceFromScores creates a new OCRConfidence instance from a slice of confidence scores.
// If the scores slice is empty, it returns nil.
func NewOCRConfidenceFromScores(scores []float32) *OCRConfidence {
	if len(scores) == 0 {
		return nil
	}

	return &OCRConfidence{
		mean: internal.Mean(scores),
		max:  slices.Max(scores),
		min:  slices.Min(scores),
	}
}

// Mean returns the mean (average) OCR confidence score.
func (c *OCRConfidence) Mean() float32 {
	return c.mean
}

// Max returns the maximum OCR confidence score.
func (c *OCRConfidence) Max() float32 {
	return c.max
}

// Min returns the minimum OCR confidence score.
func (c *OCRConfidence) Min() float32 {
	return c.min
}
