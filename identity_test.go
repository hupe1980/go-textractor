package textractor

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalyzeIDOutputParser(t *testing.T) {
	output := &textract.AnalyzeIDOutput{
		AnalyzeIDModelVersion: aws.String("1.0"),
		IdentityDocuments:     []types.IdentityDocument{},
	}

	parser := NewAnalyzeIDOutputParser(output)

	assert.NotNil(t, parser)
	assert.Equal(t, "1.0", parser.ModelVersion())
	assert.Len(t, parser.Documents(), len(output.IdentityDocuments))
}

func TestIdentityDocument(t *testing.T) {
	documentData := types.IdentityDocument{
		IdentityDocumentFields: []types.IdentityDocumentField{},
		Blocks:                 []types.Block{},
	}

	document := NewIdentityDocument(documentData)

	assert.NotNil(t, document)
	assert.Len(t, document.Fields(), len(documentData.IdentityDocumentFields))
	assert.Len(t, document.Pages(), len(documentData.Blocks))
	assert.Equal(t, IdentityDocumentTypeOther, document.Type())
}

func TestIdentityDocumentField(t *testing.T) {
	fieldData := types.IdentityDocumentField{
		Type: &types.AnalyzeIDDetections{
			Text: aws.String("FIRST_NAME"),
		},
		ValueDetection: &types.AnalyzeIDDetections{
			Text:       aws.String("John"),
			Confidence: aws.Float32(99.9),
		},
	}

	field := NewIdentityDocumentField(fieldData)

	assert.NotNil(t, field)
	assert.Equal(t, IdentityDocumentFieldType("FIRST_NAME"), field.Type())
	assert.Equal(t, "John", field.Value())
	assert.Equal(t, float32(99.9), field.Confidence())
}

func TestNormalizedIdentityDocumentFielValue(t *testing.T) {
	normalizedValueData := &types.NormalizedValue{
		Value:     aws.String("2023-01-01T12:00:00"),
		ValueType: types.ValueTypeDate,
	}

	normalizedValue := &NormalizedIdentityDocumentFielValue{
		value: normalizedValueData,
	}

	assert.NotNil(t, normalizedValue)
	assert.Equal(t, types.ValueTypeDate, normalizedValue.Type())
	assert.Equal(t, "2023-01-01T12:00:00", normalizedValue.Value())

	// Parse date value
	dateValue, err := normalizedValue.DateValue()
	assert.NoError(t, err)

	expectedDate, err := time.Parse("2006-01-02T15:04:05", "2023-01-01T12:00:00")
	assert.NoError(t, err)
	assert.Equal(t, expectedDate, dateValue)
}
