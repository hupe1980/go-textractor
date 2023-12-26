package textractor

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type Document struct {
	blockMap map[string]types.Block
	pages    []*Page
}

func NewDocument(responsePages ...*ResponsePage) *Document {
	doc := &Document{
		blockMap: make(map[string]types.Block),
	}

	var (
		currentPageBlock   *types.Block
		currentPageContent []types.Block
	)

	for _, p := range responsePages {
		for i, b := range p.Blocks {
			doc.blockMap[aws.ToString(b.Id)] = b

			if b.BlockType == types.BlockTypePage {
				if currentPageBlock != nil {
					doc.pages = append(doc.pages, NewPage(*currentPageBlock, currentPageContent, doc.blockMap))
				}

				currentPageBlock = &p.Blocks[i]
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

func (doc *Document) BlockByID(id string) *types.Block {
	b, ok := doc.blockMap[id]
	if !ok {
		return nil
	}

	return &b
}

func (doc *Document) PageNumber(n int) (*Page, error) {
	if n < 1 || n > len(doc.pages) {
		return nil, fmt.Errorf("number %d must be between 1 and %d", n, len(doc.pages))
	}

	return doc.pages[n-1], nil
}

func (doc *Document) Pages() []*Page {
	return doc.pages
}

func (doc *Document) PageCount() int {
	return len(doc.pages)
}

// Page represents a page in the document.
type Page struct {
	block  types.Block
	blocks []types.Block
	lines  []*Line
	form   *Form
	tables []*Table
}

// NewPage creates a new Page instance.
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
		default: // TODO logging?
		}
	}

	return page
}

func (p *Page) ID() string {
	return aws.ToString(p.block.Id)
}

func (p *Page) Blocks() []types.Block {
	return p.blocks
}

func (p *Page) Geometry() *Geometry {
	return NewGeometry(p.block.Geometry)
}

func (p *Page) Text() string {
	texts := make([]string, len(p.lines))
	for i, l := range p.lines {
		texts[i] = l.Text()
	}

	return strings.Join(texts, "\n")
}

func (p *Page) Form() *Form {
	return p.form
}

func (p *Page) TableCount() int {
	return len(p.tables)
}

func (p *Page) Tables() []*Table {
	return p.tables
}

func (p *Page) TableAtIndex(i int) *Table {
	return p.tables[i]
}

func (p *Page) LineCount() int {
	return len(p.lines)
}

func (p *Page) Lines() []*Line {
	return p.lines
}

func (p *Page) LineAtIndex(i int) (*Line, error) {
	if i < 0 || i >= len(p.lines) {
		return nil, fmt.Errorf("index %d must be > 0 and < %d", i, len(p.lines))
	}

	return p.lines[i], nil
}
