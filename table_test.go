package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	t.Run("Parse doc with tables", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(td.Blocks)
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, 1, len(doc.Pages()[0].Tables()))
	})

	t.Run("Table dimensions", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(td.Blocks)
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, 1, len(doc.Pages()[0].Tables()))

		table := doc.PageNumber(1).TableAtIndex(0)
		assert.Equal(t, 1, table.RowCount())

		// TODO
	})
}
