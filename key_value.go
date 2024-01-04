package textractor

import (
	"fmt"
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

// Compile time check to ensure KeyValue satisfies the LayoutChild interface.
var _ LayoutChild = (*KeyValue)(nil)

type KeyValue struct {
	base
	key   *Key
	value *Value
	page  *Page
}

func (kv *KeyValue) Key() *Key {
	return kv.key
}

func (kv *KeyValue) Value() *Value {
	return kv.value
}

// Confidence calculates the confidence score for a key value.
func (kv *KeyValue) Confidence() float64 {
	scores := make([]float64, 0)

	if kv.Key() != nil {
		scores = append(scores, kv.Key().Confidence())
	}

	if kv.Value() != nil {
		scores = append(scores, kv.Value().Confidence())
	}

	return internal.Mean(scores)
}

func (kv *KeyValue) BoundingBox() *BoundingBox {
	return NewEnclosingBoundingBox[BoundingBoxAccessor](kv.Key(), kv.Value())
}

func (kv *KeyValue) Polygon() Polygon {
	// TODO
	panic("not implemented")
}

func (kv *KeyValue) Words() []*Word {
	return internal.Concatenate(kv.Key().Words(), kv.Value().Words())
}

func (kv *KeyValue) Text(optFns ...func(*TextLinearizationOptions)) string {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	keyText := kv.Key().Text()
	keyText = fmt.Sprintf("%s%s%s", opts.KeyPrefix, keyText, opts.KeySuffix)

	valueText := kv.Value().Text()
	valueText = fmt.Sprintf("%s%s%s", opts.ValuePrefix, valueText, opts.ValueSuffix)

	if len(keyText) == 0 && len(valueText) == 0 {
		return ""
	}

	text := fmt.Sprintf("%s%s%s", keyText, opts.SameParagraphSeparator, valueText)
	if kv.Value().SelectionElement() != nil {
		text = fmt.Sprintf("%s%s%s", valueText, opts.SameParagraphSeparator, keyText)
	}

	return fmt.Sprintf("%s%s%s", opts.KeyValuePrefix, text, opts.KeyValueSuffix)
}

func (kv *KeyValue) String() string {
	if kv.Value().SelectionElement() != nil {
		return fmt.Sprintf("%s %s", kv.Value(), kv.Key())
	}

	return fmt.Sprintf("%s : %s", kv.Key(), kv.Value())
}

type Key struct {
	base
	words []*Word
}

func (k *Key) Words() []*Word {
	return k.words
}

func (k *Key) Text() string {
	texts := make([]string, len(k.words))
	for i, w := range k.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

// String returns the string representation of the key.
func (k *Key) String() string {
	return k.Text()
}

type Value struct {
	base
	words            []*Word
	selectionElement *SelectionElement
}

func (v *Value) Words() []*Word {
	return v.words
}

// SelectionElement returns the selection element associated with the table cell.
func (v *Value) SelectionElement() *SelectionElement {
	return v.selectionElement
}

func (v *Value) Text(optFns ...func(*TextLinearizationOptions)) string {
	if v.selectionElement != nil {
		return v.selectionElement.Text(optFns...)
	}

	texts := make([]string, len(v.words))
	for i, w := range v.words {
		texts[i] = w.Text()
	}

	text := strings.Join(texts, " ")

	// Replace all occurrences of \n with a space
	text = strings.ReplaceAll(text, "\n", " ")

	// Replace consecutive spaces with a single space
	for strings.Contains(text, "  ") {
		text = strings.ReplaceAll(text, "  ", " ")
	}

	return text
}

// String returns the string representation of the value.
func (v *Value) String() string {
	return v.Text()
}
