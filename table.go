package textractor

import (
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

type Table struct {
	base
	title   *TableTitle
	footers []*TableFooter
	cells   []*TableCell
}

func (t *Table) Words() []*Word {
	words := make([][]*Word, 0, len(t.cells))

	for _, c := range t.cells {
		words = append(words, c.Words())
	}

	return internal.Concatenate(words...)
}

func (t *Table) TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word) {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	words := t.Words()
	texts := []string{}

	for _, r := range t.Rows() {
		cellText := ""

		for i, c := range r.Cells() {
			if i == 0 {
				cellText += c.Text()
			} else {
				cellText += "\t" + c.Text()
			}
		}

		texts = append(texts, cellText)
	}

	return strings.Join(texts, "\n"), words
}

func (t *Table) Rows() []*TableRow {
	cellsPerRow := make(map[int][]*TableCell, 0)
	for _, c := range t.cells {
		cellsPerRow[c.rowIndex] = append(cellsPerRow[c.rowIndex], c)
	}

	rows := make([]*TableRow, len(cellsPerRow))
	for k, v := range cellsPerRow {
		rows[k-1] = &TableRow{
			cells: v,
		}
	}

	return rows
}

type TableRow struct {
	cells []*TableCell
}

func (tr *TableRow) Cells() []*TableCell {
	return tr.cells
}
