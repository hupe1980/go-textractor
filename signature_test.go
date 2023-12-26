package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestSignature(t *testing.T) {
	// Create a sample Block for testing
	sampleBlock := types.Block{
		Id:         aws.String("sampleId"),
		BlockType:  types.BlockTypeWord,
		Confidence: aws.Float32(0.95),
		Geometry: &types.Geometry{
			BoundingBox: &types.BoundingBox{
				Width:  10,
				Height: 5,
				Left:   0,
				Top:    0,
			},
		},
	}

	// Create a Signature instance for testing
	signature := NewSignature(sampleBlock)

	// Test methods of the Signature type
	assert.Equal(t, "sampleId", signature.ID())
	assert.Equal(t, float32(0.95), signature.Confidence())

	// Test Geometry method
	geometry := signature.Geometry()

	assert.NotNil(t, geometry)
	assert.Equal(t, float32(10), geometry.BoundingBox().Width())
	assert.Equal(t, float32(5), geometry.BoundingBox().Height())
	assert.Equal(t, float32(0), geometry.BoundingBox().Left())
	assert.Equal(t, float32(0), geometry.BoundingBox().Top())
}
