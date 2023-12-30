package textractor

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type IdentityDocumentField struct {
	fieldType           IdentityDocumentFieldType
	fieldTypeConfidence float32
	value               string
	valueConfidence     float32
	normalizedValue     *NormalizedIdentityDocumentFieldValue
	raw                 types.IdentityDocumentField
}

func (idf *IdentityDocumentField) FieldType() IdentityDocumentFieldType {
	return idf.fieldType
}

func (idf *IdentityDocumentField) Value() string {
	return idf.value
}

func (idf *IdentityDocumentField) IsNormalized() bool {
	return idf.normalizedValue != nil
}

func (idf *IdentityDocumentField) NormalizedValue() *NormalizedIdentityDocumentFieldValue {
	return idf.normalizedValue
}

type NormalizedIdentityDocumentFieldValue struct {
	valueType types.ValueType
	value     string
}

func (nidfv NormalizedIdentityDocumentFieldValue) ValueType() types.ValueType {
	return nidfv.valueType
}

func (nidfv NormalizedIdentityDocumentFieldValue) Value() string {
	return nidfv.value
}

func (nidfv NormalizedIdentityDocumentFieldValue) DateValue() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", nidfv.value)
}
