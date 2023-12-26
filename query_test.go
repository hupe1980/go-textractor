package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestQueryResult(t *testing.T) {
	block := types.Block{
		Id:         aws.String("123"),
		BlockType:  types.BlockTypeQueryResult,
		Confidence: aws.Float32(99.123),
		Text:       aws.String("Sample Text"),
		Geometry: &types.Geometry{
			BoundingBox: &types.BoundingBox{
				Width:  10,
				Height: 5,
				Left:   0,
				Top:    0,
			},
		},
	}

	queryResult := NewQueryResult(block)

	assert.Equal(t, "123", queryResult.ID())
	assert.Equal(t, float32(99.123), queryResult.Confidence())
	assert.Equal(t, "Sample Text", queryResult.Text())

	geometry := queryResult.Geometry()

	assert.NotNil(t, geometry)
	assert.Equal(t, float32(10), geometry.BoundingBox().Width())
	assert.Equal(t, float32(5), geometry.BoundingBox().Height())
	assert.Equal(t, float32(0), geometry.BoundingBox().Left())
	assert.Equal(t, float32(0), geometry.BoundingBox().Top())
}

func TestQuery(t *testing.T) {
	blockMap := map[string]types.Block{
		"1": {
			Id:         aws.String("1"),
			BlockType:  types.BlockTypeQueryResult,
			Text:       aws.String("foo"),
			Confidence: aws.Float32(80),
		},
		"2": {
			Id:         aws.String("2"),
			BlockType:  types.BlockTypeQueryResult,
			Text:       aws.String("bar"),
			Confidence: aws.Float32(90),
		},
	}

	query := NewQuery(types.Block{
		Query: &types.Query{
			Text:  aws.String("Query Text"),
			Alias: aws.String("Query Alias"),
		},
		Relationships: []types.Relationship{{
			Ids:  []string{"1", "2"},
			Type: types.RelationshipTypeAnswer,
		}},
	}, blockMap)

	assert.Equal(t, "Query Alias", query.Alias())
	assert.Equal(t, "Query Text", query.Text())

	t.Run("ResultsByConfidence", func(t *testing.T) {
		results := query.ResultsByConfidence()
		assert.Len(t, results, 2)
		assert.Equal(t, "bar", results[0].Text())
		assert.Equal(t, "foo", results[1].Text())
	})

	t.Run("TopResult", func(t *testing.T) {
		topResult := query.TopResult()
		assert.NotNil(t, topResult)
		assert.Equal(t, "bar", topResult.Text())
	})
}
