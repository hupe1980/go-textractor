package textractor

import (
	"testing"

	"github.com/hupe1980/go-textractor/internal"
	"github.com/stretchr/testify/assert"
)

func TestForm(t *testing.T) {
	t.Run("Parse doc with forms", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(&ResponsePage{Blocks: td.Blocks})
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, 4, len(doc.Pages()[0].Form().Fields()))

		f1 := doc.Pages()[0].Form().FieldByKey("Home Address:")
		assert.Equal(t, "123 Any Street. Any Town. USA", f1.Value().Text())
		assert.Equal(t, "03c0345c-f42d-4bea-864d-60a8d1e890fb", f1.Key().ID())
		assert.Equal(t, "1059d4c4-dd84-4995-be8f-24c0e8b12a6a", f1.Value().ID())
		assert.Equal(t, float32(65.03003), f1.Key().Confidence())
		assert.Equal(t, float32(65.03003), f1.Value().Confidence())
		assert.Equal(t, internal.Mean([]float32{f1.Key().Confidence(), f1.Value().Confidence()}), f1.Confidence())

		f2 := doc.Pages()[0].Form().FieldByKey("Mailing Address:")
		assert.Equal(t, "same as home address", f2.Value().Text())
		assert.Equal(t, "b7684bc3-cec7-4f4e-a2bc-ac0866094ac6", f2.Key().ID())
		assert.Equal(t, "c327f8af-fd9f-47d8-bf47-b8fe34719fd3", f2.Value().ID())
		assert.Equal(t, float32(61.291843), f2.Key().Confidence())
		assert.Equal(t, float32(61.291843), f2.Value().Confidence())
		assert.Equal(t, internal.Mean([]float32{f2.Key().Confidence(), f2.Value().Confidence()}), f2.Confidence())

		f3 := doc.Pages()[0].Form().FieldByKey("Phone Number:")
		assert.Equal(t, "555-0100", f3.Value().Text())
		assert.Equal(t, "dd8dbecf-73f1-49e2-bc9f-7169a133f6dd", f3.Key().ID())
		assert.Equal(t, "a8acd770-2d5a-4799-9e0d-17a96b6da85a", f3.Value().ID())

		f4 := doc.Pages()[0].Form().FieldByKey("Full Name:")
		assert.Equal(t, "Jane Doe", f4.Value().Text())
		assert.Equal(t, "100c9244-9c74-4166-82b5-1e9890cf455d", f4.Key().ID())
		assert.Equal(t, "3f728c36-56f2-487c-a3e2-3cafe49f7da9", f4.Value().ID())
	})
}
