package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	t.Run("Parse doc with pages", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(&ResponsePage{Blocks: td.Blocks})
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, len(doc.Pages()), doc.PageCount())

		firstPage, err := doc.PageNumber(1)
		assert.NoError(t, err)
		assert.Equal(t, doc.Pages()[0], firstPage)
	})
}
