package textractor

import (
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor/internal"
)

// Form represents a form extracted from a document.
type Form struct {
	fieldMap map[string]*Field
}

// NewForm creates a new Form instance.
func NewForm() *Form {
	return &Form{
		fieldMap: make(map[string]*Field),
	}
}

// AddField adds a field to the form, replacing it if a field with the
// same key already and lower confidence exists.
func (f *Form) AddField(field *Field) {
	ef, ok := f.fieldMap[field.Key().String()]
	if ok && ef.Confidence() > field.Confidence() {
		return
	}

	key := strings.TrimSpace(field.Key().String())

	f.fieldMap[key] = field
}

// FieldByKey retrieves a field from the form by its key.
func (f *Form) FieldByKey(key string) *Field {
	field, ok := f.fieldMap[key]
	if !ok {
		return nil
	}

	return field
}

// SearchFieldByKey searches for fields in the form with a key containing the specified string.
// It performs a case-insensitive search on the key text.
func (f *Form) SearchFieldByKey(key string) []*Field {
	searchKey := strings.ToLower(key)

	var result []*Field

	for _, field := range f.Fields() {
		if key := field.Key(); key != nil {
			if strings.Contains(strings.ToLower(key.Text()), searchKey) {
				result = append(result, field)
			}
		}
	}

	return result
}

// Fields returns all fields in the form.
func (f *Form) Fields() []*Field {
	return internal.Values(f.fieldMap)
}

// FieldKey represents the key part of a form field.
type FieldKey struct {
	words []*Word
	content
}

// NewFieldKey creates a new FieldKey instance.
func NewFieldKey(block types.Block, ids []string, blockMap map[string]types.Block) *FieldKey {
	k := &FieldKey{
		content: content{block},
	}

	for _, i := range ids {
		b := blockMap[i]
		if b.BlockType == types.BlockTypeWord {
			k.words = append(k.words, NewWord(b))
		}
	}

	return k
}

// Text returns the text representation of the field key.
func (fk *FieldKey) Text() string {
	texts := make([]string, len(fk.words))
	for i, w := range fk.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

// Words returns the words constituting the field key.
func (fk *FieldKey) Words() []*Word {
	return fk.words
}

// OCRConfidence calculates the OCR confidence for the field key.
func (fk *FieldKey) OCRConfidence() *OCRConfidence {
	scores := make([]float32, len(fk.words))
	for i, w := range fk.words {
		scores[i] = w.Confidence()
	}

	return NewOCRConfidenceFromScores(scores)
}

// String returns the string representation of the field key.
func (fk *FieldKey) String() string {
	return fk.Text()
}

// FieldValue represents the value part of a form field.
type FieldValue struct {
	words            []*Word
	selectionElement *SelectionElement
	content
}

// NewFieldValue creates a new FieldValue instance.
func NewFieldValue(block types.Block, ids []string, blockMap map[string]types.Block) *FieldValue {
	v := &FieldValue{
		content: content{block},
	}

	for _, i := range ids {
		b := blockMap[i]
		if b.BlockType == types.BlockTypeWord {
			v.words = append(v.words, NewWord(b))
		} else if b.BlockType == types.BlockTypeSelectionElement {
			v.selectionElement = NewSelectionElement(b)
		}
	}

	return v
}

// Text returns the text representation of the field value.
func (fv *FieldValue) Text() string {
	texts := make([]string, len(fv.words))
	for i, w := range fv.words {
		texts[i] = w.Text()
	}

	if fv.selectionElement != nil {
		texts = append(texts, string(fv.selectionElement.Status()))
	}

	return strings.Join(texts, " ")
}

// Words returns the words constituting the field value.
func (fv *FieldValue) Words() []*Word {
	return fv.words
}

// SelectionElement returns the selection element associated with the field value.
func (fv *FieldValue) SelectionElement() *SelectionElement {
	return fv.selectionElement
}

// OCRConfidence calculates the OCR confidence for the field value.
func (fv *FieldValue) OCRConfidence() *OCRConfidence {
	scores := make([]float32, len(fv.words))
	for i, w := range fv.words {
		scores[i] = w.Confidence()
	}

	if fv.selectionElement != nil {
		scores = append(scores, fv.selectionElement.Confidence())
	}

	return NewOCRConfidenceFromScores(scores)
}

// String returns the string representation of the field value.
func (fv *FieldValue) String() string {
	return fv.Text()
}

// Field represents a form field, consisting of a key and a value.
type Field struct {
	key   *FieldKey
	value *FieldValue
}

// NewField creates a new Field instance.
func NewField(block types.Block, blockMap map[string]types.Block) *Field {
	field := &Field{}

	for _, r := range block.Relationships {
		if r.Type == types.RelationshipTypeChild {
			field.key = NewFieldKey(block, r.Ids, blockMap)
		} else if r.Type == types.RelationshipTypeValue {
		valueLoop:
			for _, i := range r.Ids {
				v := blockMap[i]
				if slices.Contains(v.EntityTypes, types.EntityTypeValue) {
					for _, vr := range v.Relationships {
						if vr.Type == types.RelationshipTypeChild {
							field.value = NewFieldValue(v, vr.Ids, blockMap)
							break valueLoop
						}
					}
				}
			}
		}
	}

	return field
}

// Confidence calculates the confidence score for the form field.
func (f *Field) Confidence() float32 {
	scores := []float32{}

	if f.Key() != nil {
		scores = append(scores, f.Key().Confidence())
	}

	if f.Value() != nil {
		scores = append(scores, f.Value().Confidence())
	}

	return internal.Mean(scores)
}

// OCRConfidence calculates the OCR confidence for the form field.
func (f *Field) OCRConfidence() *OCRConfidence {
	scores := make([]float32, 0)

	if f.Key() != nil {
		for _, w := range f.Key().Words() {
			scores = append(scores, w.Confidence())
		}
	}

	if f.Value() != nil {
		for _, w := range f.Value().Words() {
			scores = append(scores, w.Confidence())
		}

		if f.Value().SelectionElement() != nil {
			scores = append(scores, f.Value().SelectionElement().Confidence())
		}
	}

	return NewOCRConfidenceFromScores(scores)
}

// Key returns the key part of the form field.
func (f *Field) Key() *FieldKey {
	return f.key
}

// Value returns the value part of the form field.
func (f *Field) Value() *FieldValue {
	return f.value
}
