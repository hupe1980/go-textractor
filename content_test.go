package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	t.Run("Parse doc with liens and words", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(&ResponsePage{Blocks: td.Blocks})
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, 22, len(doc.Pages()[0].Lines()))

		wc := 0
		for _, l := range doc.Pages()[0].Lines() {
			wc += len(l.Words())
		}
		assert.Equal(t, 53, wc)
	})
}
