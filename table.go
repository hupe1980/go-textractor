package textractor

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type Cell struct {
	words            []*Word
	selectionElement *SelectionElement
	content
}

func NewCell(block types.Block, blockMap map[string]types.Block) *Cell {
	cell := &Cell{
		content: content{block},
	}

	for _, r := range block.Relationships {
		if r.Type == types.RelationshipTypeChild {
			for _, i := range r.Ids {
				b := blockMap[i]
				if b.BlockType == types.BlockTypeWord {
					cell.words = append(cell.words, NewWord(b))
				} else if b.BlockType == types.BlockTypeSelectionElement {
					cell.selectionElement = NewSelectionElement(b)
				}
			}
		}
	}

	return cell
}

func (c *Cell) ColumnIndex() int32 {
	return aws.ToInt32(c.Block().ColumnIndex)
}

func (c *Cell) ColumnSpan() int32 {
	if c.Block().ColumnSpan != nil {
		return aws.ToInt32(c.Block().ColumnSpan)
	}

	return 1
}

func (c *Cell) RowIndex() int32 {
	return aws.ToInt32(c.Block().RowIndex)
}

func (c *Cell) RowSpan() int32 {
	if c.Block().RowSpan != nil {
		return aws.ToInt32(c.Block().RowSpan)
	}

	return 1
}

func (c *Cell) Text() string {
	texts := make([]string, len(c.words))
	for i, w := range c.words {
		texts[i] = w.Text()
	}

	if c.selectionElement != nil {
		texts = append(texts, string(c.selectionElement.Status()))
	}

	return strings.Join(texts, " ")
}

type Row struct {
	cells []*Cell
	*content
}

func NewRow() *Row {
	return &Row{}
}

func (r *Row) AddCell(cell *Cell) {
	r.cells = append(r.cells, cell)
}

func (r *Row) CellCount() int {
	return len(r.cells)
}

func (r *Row) CellAt(i int) *Cell {
	return r.cells[i]
}

func (r *Row) Cells() []*Cell {
	return r.cells
}

type Table struct {
	block types.Block
	rows  []*Row
	*content
}

func NewTable(block types.Block, blockMap map[string]types.Block) *Table {
	table := &Table{
		block:   block,
		content: &content{block},
	}

	ri := int32(1)
	row := NewRow()

	for _, r := range block.Relationships {
		if r.Type == types.RelationshipTypeChild {
			for _, i := range r.Ids {
				cell := NewCell(blockMap[i], blockMap)
				if cell.RowIndex() > ri {
					table.rows = append(table.rows, row)
					row = NewRow()
					ri = cell.RowIndex()
				}

				row.AddCell(cell)
			}

			if row != nil && row.CellCount() > 0 {
				table.rows = append(table.rows, row)
			}
		}
	}

	return table
}

func (t *Table) RowCount() int {
	return len(t.rows)
}

func (t *Table) Rows() []*Row {
	return t.rows
}

func (t *Table) RowAt(rowIndex int) *Row {
	return t.rows[rowIndex]
}

func (t *Table) CellAt(rowIndex, columnIndex int) *Cell {
	row := t.RowAt(rowIndex)
	return row.CellAt(columnIndex)
}
