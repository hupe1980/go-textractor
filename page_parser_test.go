package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestPageParser(t *testing.T) {
	t.Run("createSignatures", func(t *testing.T) {
		signatureBlock1 := types.Block{
			Id:         aws.String("id1"),
			BlockType:  types.BlockTypeSignature,
			Confidence: aws.Float32(0.95),
			Geometry: &types.Geometry{
				BoundingBox: &types.BoundingBox{},
			},
		}

		signatureBlock2 := types.Block{
			Id:         aws.String("id2"),
			BlockType:  types.BlockTypeSignature,
			Confidence: aws.Float32(0.98),
			Geometry: &types.Geometry{
				BoundingBox: &types.BoundingBox{},
			},
		}

		bp := newBlockParser([]types.Block{signatureBlock1, signatureBlock2})

		page := &Page{
			childIDs: []string{"id1", "id2"},
		}

		pp := newPageParser(bp, page)

		resultSignatures := pp.createSignatures()

		assert.Len(t, resultSignatures, 2)
		assert.Equal(t, "id1", resultSignatures[0].id)
		assert.Equal(t, "id2", resultSignatures[1].id)
	})
}
