package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestBase(t *testing.T) {
	t.Run("IDMethod", func(t *testing.T) {
		b := base{id: "123"}
		assert.Equal(t, "123", b.ID(), "ID method does not return expected value")
	})

	t.Run("ConfidenceMethod", func(t *testing.T) {
		b := base{confidence: 0.9}
		assert.Equal(t, 0.9, b.Confidence(), "Confidence method does not return expected value")
	})

	t.Run("BlockTypeMethod", func(t *testing.T) {
		b := base{blockType: types.BlockTypePage}
		assert.Equal(t, types.BlockTypePage, b.BlockType(), "BlockType method does not return expected value")
	})

	t.Run("BoundingBoxMethod", func(t *testing.T) {
		boundingBox := &BoundingBox{left: 1, top: 2, width: 3, height: 4}
		b := base{boundingBox: boundingBox}
		assert.Equal(t, boundingBox, b.BoundingBox(), "BoundingBox method does not return expected value")
	})

	t.Run("PolygonMethod", func(t *testing.T) {
		polygon := Polygon{{1, 2}, {3, 4}}
		b := base{polygon: polygon}
		assert.Equal(t, polygon, b.Polygon(), "Polygon method does not return expected value")
	})

	t.Run("PageNumberMethod", func(t *testing.T) {
		page := &Page{number: 5}
		b := base{page: page}
		assert.Equal(t, 5, b.PageNumber(), "PageNumber method does not return expected value")
	})

	t.Run("RawMethod", func(t *testing.T) {
		rawBlock := types.Block{Id: aws.String("456"), BlockType: types.BlockTypeWord}
		b := base{raw: rawBlock}
		assert.Equal(t, rawBlock, b.Raw(), "Raw method does not return expected value")
	})
}
