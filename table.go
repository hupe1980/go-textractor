package textractor

type Table struct {
	base
	title   *TableTitle
	footers []*TableFooter
	cells   []*TableCell
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
