package textractor

import (
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor/internal"
)

type Form struct {
	fieldMap map[string]*Field
}

func NewForm() *Form {
	return &Form{
		fieldMap: make(map[string]*Field),
	}
}

func (f *Form) AddField(field *Field) {
	ef, ok := f.fieldMap[field.Key().String()]
	if ok && ef.Confidence() > field.Confidence() {
		return
	}

	key := strings.TrimSpace(field.Key().String())

	f.fieldMap[key] = field
}

func (f *Form) FieldByKey(key string) *Field {
	field, ok := f.fieldMap[key]
	if !ok {
		return nil
	}

	return field
}

func (f *Form) Fields() []*Field {
	return internal.Values(f.fieldMap)
}

type FieldKey struct {
	words []*Word
	content
}

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

func (fk *FieldKey) Text() string {
	texts := make([]string, len(fk.words))
	for i, w := range fk.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

func (fk *FieldKey) String() string {
	return fk.Text()
}

type FieldValue struct {
	words            []*Word
	selectionElement *SelectionElement
	content
}

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

func (fv *FieldValue) String() string {
	return fv.Text()
}

type Field struct {
	key   *FieldKey
	value *FieldValue
}

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

func (f *Field) Confidence() float32 {
	scores := []float32{}

	if f.key != nil {
		scores = append(scores, f.key.Confidence())
	}

	if f.value != nil {
		scores = append(scores, f.value.Confidence())
	}

	return internal.Mean(scores)
}

func (f *Field) Key() *FieldKey {
	return f.key
}

func (f *Field) Value() *FieldValue {
	return f.value
}
