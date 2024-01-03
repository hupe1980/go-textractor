package textractor

import (
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor/internal"
)

type Cell interface {
	Words() []*Word
	Text(optFns ...func(*TextLinearizationOptions)) string
	Confidence() float64
}

type TableMergedCell struct {
	cell
	cells []*TableCell
}

func (tmc *TableMergedCell) Words() []*Word {
	words := make([][]*Word, 0, len(tmc.cells))

	for _, c := range tmc.cells {
		words = append(words, c.Words())
	}

	return internal.Concatenate(words...)
}

func (tmc *TableMergedCell) Text(_ ...func(*TextLinearizationOptions)) string {
	words := tmc.Words()

	texts := make([]string, len(words))
	for i, w := range words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

type TableCell struct {
	cell
	words            []*Word
	selectionElement *SelectionElement
}

func (tc *TableCell) Words() []*Word {
	return tc.words
}

func (tc *TableCell) Text(optFns ...func(*TextLinearizationOptions)) string {
	if tc.selectionElement != nil {
		return tc.selectionElement.Text(optFns...)
	}

	texts := make([]string, len(tc.words))
	for i, w := range tc.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

type cell struct {
	base
	rowIndex    int
	columnIndex int
	rowSpan     int
	columnSpan  int
	entityTypes []types.EntityType
}

func newCell(cb types.Block, p *Page) cell {
	return cell{
		base:        newBase(cb, p),
		rowIndex:    int(aws.ToInt32(cb.RowIndex)),
		columnIndex: int(aws.ToInt32(cb.ColumnIndex)),
		rowSpan:     int(aws.ToInt32(cb.RowSpan)),
		columnSpan:  int(aws.ToInt32(cb.ColumnSpan)),
	}
}

func (c *cell) IsColumnHeader() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeColumnHeader)
}

func (c *cell) IsTableTitle() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableTitle)
}

func (c *cell) IsTableFooter() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableFooter)
}

func (c *cell) IsTableSummary() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableSummary)
}

func (c *cell) IsTableSectionTitle() bool {
	return slices.Contains(c.entityTypes, types.EntityTypeTableSectionTitle)
}
