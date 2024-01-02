package textractor

import (
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type TableCell struct {
	base
	rowIndex    int
	columnIndex int
	rowSpan     int
	columnSpan  int
	entityTypes []types.EntityType
	words       []*Word
}

func (tc *TableCell) Words() []*Word {
	return tc.words
}

func (tc *TableCell) Text() string {
	texts := make([]string, len(tc.words))
	for i, w := range tc.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

func (tc *TableCell) IsColumnHeader() bool {
	return slices.Contains(tc.entityTypes, types.EntityTypeColumnHeader)
}

func (tc *TableCell) IsTableTitle() bool {
	return slices.Contains(tc.entityTypes, types.EntityTypeTableTitle)
}

func (tc *TableCell) IsTableFooter() bool {
	return slices.Contains(tc.entityTypes, types.EntityTypeTableFooter)
}

func (tc *TableCell) IsTableSummary() bool {
	return slices.Contains(tc.entityTypes, types.EntityTypeTableSummary)
}

func (tc *TableCell) IsTableSectionTitle() bool {
	return slices.Contains(tc.entityTypes, types.EntityTypeTableSectionTitle)
}
