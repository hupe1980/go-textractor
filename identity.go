package textractor

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// AnalyzeIDOutputParser is a parser for the output of the AnalyzeID operation in Amazon Textract.
type AnalyzeIDOutputParser struct {
	analyzeIDModelVersion string
	documents             []*IdentityDocument
}

// NewAnalyzeIDOutputParser creates a new AnalyzeIDOutputParser instance from the output of the AnalyzeID operation.
func NewAnalyzeIDOutputParser(output *textract.AnalyzeIDOutput) *AnalyzeIDOutputParser {
	docs := make([]*IdentityDocument, len(output.IdentityDocuments))

	for i, d := range output.IdentityDocuments {
		docs[i] = NewIdentityDocument(d)
	}

	return &AnalyzeIDOutputParser{
		analyzeIDModelVersion: aws.ToString(output.AnalyzeIDModelVersion),
		documents:             docs,
	}
}

// ModelVersion returns the model version used for the analysis.
func (p *AnalyzeIDOutputParser) ModelVersion() string {
	return p.analyzeIDModelVersion
}

// Documents returns the identity documents extracted from the analysis.
func (p *AnalyzeIDOutputParser) Documents() []*IdentityDocument {
	return p.documents
}

// IdentityDocumentFieldType represents the type of fields in an identity document.
type IdentityDocumentFieldType string

const (
	IdentityDocumentFieldTypeFirstName        IdentityDocumentFieldType = "FIRST_NAME"
	IdentityDocumentFieldTypeLastName         IdentityDocumentFieldType = "LAST_NAME"
	IdentityDocumentFieldTypeMiddleName       IdentityDocumentFieldType = "MIDDLE_NAME"
	IdentityDocumentFieldTypeSuffix           IdentityDocumentFieldType = "Suffix"
	IdentityDocumentFieldTypeCityInAddress    IdentityDocumentFieldType = "CITY_IN_ADDRESS"
	IdentityDocumentFieldTypeZipCodeInAddress IdentityDocumentFieldType = "ZIP_CODE_IN_ADDRESS"
	IdentityDocumentFieldTypeStateInAddress   IdentityDocumentFieldType = "STATE_IN_ADDRESS"
	IdentityDocumentFieldTypeStateName        IdentityDocumentFieldType = "STATE_NAME"
	IdentityDocumentFieldTypeDocumentNumber   IdentityDocumentFieldType = "DOCUMENT_NUMBER"
	IdentityDocumentFieldTypeExpirationDate   IdentityDocumentFieldType = "EXPIRATION_DATE"
	IdentityDocumentFieldTypeDateOfBirth      IdentityDocumentFieldType = "DATE_OF_BIRTH"
	IdentityDocumentFieldTypeDateOfIssue      IdentityDocumentFieldType = "DATE_OF_ISSUE"
	IdentityDocumentFieldTypeIDType           IdentityDocumentFieldType = "ID_TYPE"
	IdentityDocumentFieldTypeEndorsements     IdentityDocumentFieldType = "ENDORSEMENTS"
	IdentityDocumentFieldTypeVeteran          IdentityDocumentFieldType = "VETERAN"
	IdentityDocumentFieldTypeRestrictions     IdentityDocumentFieldType = "RESTRICTIONS"
	IdentityDocumentFieldTypeClass            IdentityDocumentFieldType = "CLASS"
	IdentityDocumentFieldTypeAddress          IdentityDocumentFieldType = "ADDRESS"
	IdentityDocumentFieldTypeCounty           IdentityDocumentFieldType = "COUNTY"
	IdentityDocumentFieldTypePlaceOfBirth     IdentityDocumentFieldType = "PLACE_OF_BIRTH"
	IdentityDocumentFieldTypeOther            IdentityDocumentFieldType = "Other"
)

// IdentityDocumentType represents the type of an identity document.
type IdentityDocumentType string

const (
	IdentityDocumentTypeDrivingLicense IdentityDocumentType = "DRIVER LICENSE FRONT"
	IdentityDocumentTypePassport       IdentityDocumentType = "PASSPORT"
	IdentityDocumentTypeOther          IdentityDocumentType = "OTHER"
)

// IdentityDocument represents an extracted identity document.
type IdentityDocument struct {
	document  *Document
	fields    []*IdentityDocumentField
	fieldsMap map[IdentityDocumentFieldType]*IdentityDocumentField
}

// NewIdentityDocument creates a new IdentityDocument instance from the extracted identity document in Amazon Textract.
func NewIdentityDocument(document types.IdentityDocument) *IdentityDocument {
	fields := make([]*IdentityDocumentField, len(document.IdentityDocumentFields))
	fieldsMap := make(map[IdentityDocumentFieldType]*IdentityDocumentField, len(document.IdentityDocumentFields))

	for i, f := range document.IdentityDocumentFields {
		field := NewIdentityDocumentField(f)
		fields[i] = field
		fieldsMap[field.Type()] = field
	}

	return &IdentityDocument{
		document:  NewDocument(document.Blocks),
		fields:    fields,
		fieldsMap: fieldsMap,
	}
}

// Type returns the type of the identity document.
func (doc *IdentityDocument) Type() IdentityDocumentType {
	if f := doc.FieldByType(IdentityDocumentFieldTypeIDType); f != nil {
		return IdentityDocumentType(f.Value())
	}

	return IdentityDocumentTypeOther
}

// FieldCount returns the number of fields in the identity document.
func (doc *IdentityDocument) FieldCount() int {
	return len(doc.fields)
}

// Fields returns the fields in the identity document.
func (doc *IdentityDocument) Fields() []*IdentityDocumentField {
	return doc.fields
}

// FieldByType returns the field in the identity document of the specified type.
func (doc *IdentityDocument) FieldByType(ft IdentityDocumentFieldType) *IdentityDocumentField {
	field, ok := doc.fieldsMap[ft]
	if !ok {
		return nil
	}

	return field
}

// Pages returns all pages in the document.
func (doc *IdentityDocument) Pages() []*Page {
	return doc.document.Pages()
}

// IdentityDocumentField represents a field within an identity document.
type IdentityDocumentField struct {
	field types.IdentityDocumentField
}

// NewIdentityDocumentField creates a new IdentityDocumentField instance from an identity document field in Amazon Textract.
func NewIdentityDocumentField(field types.IdentityDocumentField) *IdentityDocumentField {
	return &IdentityDocumentField{
		field: field,
	}
}

// Type returns the type of the identity document field.
func (idf *IdentityDocumentField) Type() IdentityDocumentFieldType {
	t := aws.ToString(idf.field.Type.Text)
	if t != "" {
		return IdentityDocumentFieldType(t)
	}

	return IdentityDocumentFieldTypeOther
}

// Value returns the value of the identity document field.
func (idf *IdentityDocumentField) Value() string {
	return aws.ToString(idf.field.ValueDetection.Text)
}

// NormalizedIdentityDocumentFielValue represents a normalized value of an identity document field.
type NormalizedIdentityDocumentFielValue struct {
	value *types.NormalizedValue
}

// Type returns the value type of the normalized value.
func (nidfv *NormalizedIdentityDocumentFielValue) Type() types.ValueType {
	return nidfv.value.ValueType
}

// Value returns the string value of the normalized value.
func (nidfv *NormalizedIdentityDocumentFielValue) Value() string {
	return aws.ToString(nidfv.value.Value)
}

// DateValue returns the time.Time value of the normalized value if it is a date.
func (nidfv *NormalizedIdentityDocumentFielValue) DateValue() (time.Time, error) {
	return time.Parse("2006-01-02T15:04:05", aws.ToString(nidfv.value.Value))
}

// NormalizedValue returns the normalized value of the identity document field if available.
func (idf *IdentityDocumentField) NormalizedValue() *NormalizedIdentityDocumentFielValue {
	if idf.IsNormalized() {
		return &NormalizedIdentityDocumentFielValue{
			value: idf.field.ValueDetection.NormalizedValue,
		}
	}

	return nil
}

// IsNormalized returns true if the identity document field has a normalized value.
func (idf *IdentityDocumentField) IsNormalized() bool {
	return idf.field.ValueDetection.NormalizedValue != nil
}

// Confidence returns the confidence score of the identity document field.
func (idf *IdentityDocumentField) Confidence() float32 {
	return aws.ToFloat32(idf.field.ValueDetection.Confidence)
}
