package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuery(t *testing.T) {
	t.Run("Text", func(t *testing.T) {
		query := Query{text: "TestText"}
		assert.Equal(t, "TestText", query.Text(), "Text method does not return expected value")
	})

	t.Run("Alias", func(t *testing.T) {
		query := Query{alias: "TestAlias"}
		assert.Equal(t, "TestAlias", query.Alias(), "Alias method does not return expected value")
	})

	t.Run("HasResult_NoResult", func(t *testing.T) {
		query := Query{}
		assert.False(t, query.HasResult())
	})

	t.Run("TopResult_NoResults", func(t *testing.T) {
		query := Query{}
		assert.Nil(t, query.TopResult(), "TopResult method should return nil for no results")
	})

	t.Run("ResultsByConfidence_NoResults", func(t *testing.T) {
		query := Query{}
		assert.Empty(t, query.ResultsByConfidence(), "ResultsByConfidence method should return an empty slice for no results")
	})

	t.Run("HasResult_WithResults", func(t *testing.T) {
		query := Query{
			results: []*QueryResult{
				{base: base{confidence: 0.8}},
				{base: base{confidence: 0.9}},
				{base: base{confidence: 0.7}},
			},
		}
		assert.True(t, query.HasResult())
	})

	t.Run("TopResult_WithResults", func(t *testing.T) {
		query := Query{
			results: []*QueryResult{
				{base: base{confidence: 0.8}},
				{base: base{confidence: 0.9}},
				{base: base{confidence: 0.7}},
			},
		}
		expectedResult := &QueryResult{base: base{confidence: 0.9}}
		assert.Equal(t, expectedResult, query.TopResult(), "TopResult method does not return the top result")
	})

	t.Run("ResultsByConfidence_WithResults", func(t *testing.T) {
		query := Query{
			results: []*QueryResult{
				{base: base{confidence: 0.8}},
				{base: base{confidence: 0.9}},
				{base: base{confidence: 0.7}},
			},
		}
		expectedResults := []*QueryResult{
			{base: base{confidence: 0.9}},
			{base: base{confidence: 0.8}},
			{base: base{confidence: 0.7}},
		}
		assert.Equal(t, expectedResults, query.ResultsByConfidence(), "ResultsByConfidence method does not return expected results")
	})
}

func TestQueryResult(t *testing.T) {
	t.Run("Text", func(t *testing.T) {
		qr := QueryResult{text: "TestText"}
		assert.Equal(t, "TestText", qr.Text(), "Text method does not return expected value")
	})
}
