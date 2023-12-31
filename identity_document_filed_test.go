package textractor

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestIdentityDocumentField(t *testing.T) {
	// Test data
	rawField := types.IdentityDocumentField{
		Type: &types.AnalyzeIDDetections{
			Text: aws.String(string(IdentityDocumentFieldTypeDateOfBirth)),
		},
		ValueDetection: &types.AnalyzeIDDetections{
			Text:       aws.String("1990-01-01T15:04:05"),
			Confidence: aws.Float32(0.95),
		},
	}

	// Create an IdentityDocumentField
	normalizedValue := &NormalizedIdentityDocumentFieldValue{
		valueType: types.ValueTypeDate,
		value:     "1990-01-01T15:04:05",
	}

	identityField := &IdentityDocumentField{
		fieldType:       IdentityDocumentFieldType(aws.ToString(rawField.Type.Text)),
		value:           aws.ToString(rawField.ValueDetection.Text),
		confidence:      aws.ToFloat32(rawField.ValueDetection.Confidence),
		normalizedValue: normalizedValue,
		raw:             rawField,
	}

	// Test methods of IdentityDocumentField
	assert.Equal(t, IdentityDocumentFieldTypeDateOfBirth, identityField.FieldType())
	assert.Equal(t, "1990-01-01T15:04:05", identityField.Value())
	assert.Equal(t, float32(0.95), identityField.Confidence())
	assert.True(t, identityField.IsNormalized())
	assert.Equal(t, normalizedValue, identityField.NormalizedValue())
}

func TestNormalizedIdentityDocumentFieldValue(t *testing.T) {
	// Test data
	rawValue := types.ValueTypeDate
	value := "1990-01-01T15:04:05"

	// Create a NormalizedIdentityDocumentFieldValue
	normalizedValue := &NormalizedIdentityDocumentFieldValue{
		valueType: rawValue,
		value:     value,
	}

	// Test methods of NormalizedIdentityDocumentFieldValue
	assert.Equal(t, rawValue, normalizedValue.ValueType())
	assert.Equal(t, value, normalizedValue.Value())

	// Test DateValue method
	dateValue, err := normalizedValue.DateValue()
	assert.NoError(t, err)

	expectedDate, err := time.Parse("2006-01-02T15:04:05", "1990-01-01T15:04:05")
	assert.NoError(t, err)
	assert.Equal(t, expectedDate, dateValue)
}
