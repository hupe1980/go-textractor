package textractor

import (
	"testing"

	"github.com/hupe1980/go-textractor/internal"
	"github.com/stretchr/testify/assert"
)

func TestDocument(t *testing.T) {
	t.Run("Parse doc with pages", func(t *testing.T) {
		td, err := loadTestdata("testdata/test-response.json")
		assert.NoError(t, err)

		doc := NewDocument(td.Blocks)
		assert.Equal(t, 1, len(doc.Pages()))
		assert.Equal(t, len(doc.Pages()), doc.PageCount())

		firstPage := doc.PageNumber(1)
		assert.Equal(t, doc.Pages()[0], firstPage)
	})
}

func TestOCRConfidence(t *testing.T) {
	t.Run("NewOCRConfidenceFromScores with empty scores", func(t *testing.T) {
		confidence := NewOCRConfidenceFromScores([]float32{})
		assert.Nil(t, confidence, "Expected nil confidence for empty scores")
	})

	t.Run("NewOCRConfidenceFromScores with non-empty scores", func(t *testing.T) {
		scores := []float32{0.8, 0.9, 0.75, 0.95}
		confidence := NewOCRConfidenceFromScores(scores)
		assert.NotNil(t, confidence, "Expected non-nil confidence for non-empty scores")

		assert.InDelta(t, internal.Mean(scores), confidence.Mean(), 0.0001, "Mean value mismatch")
		assert.Equal(t, float32(0.95), confidence.Max(), "Max value mismatch")
		assert.Equal(t, float32(0.75), confidence.Min(), "Min value mismatch")
	})
}
