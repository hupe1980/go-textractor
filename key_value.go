package textractor

import (
	"fmt"
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

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
func (kv *KeyValue) Confidence() float32 {
	scores := make([]float32, 0)

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

func (kv *KeyValue) Polygon() []*Point {
	// TODO
	panic("not implemented")
}

func (kv *KeyValue) Words() []*Word {
	return internal.Concatenate(kv.Key().Words(), kv.Value().Words())
}

func (kv *KeyValue) TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word) {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	keyText := kv.Key().Text()
	valueText := kv.Value().Text()

	if len(keyText) == 0 && len(valueText) == 0 {
		return "", nil
	}

	text := fmt.Sprintf("%s %s", keyText, valueText)

	return text, kv.value.words
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
	if v.selectionElement != nil {
		return v.selectionElement.Words()
	}

	return v.words
}

func (v *Value) Text() string {
	text, _ := v.TextAndWords()
	return text
}

func (v *Value) TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word) {
	if v.selectionElement != nil {
		return v.selectionElement.TextAndWords(optFns...)
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

	words := v.words

	return text, words
}

// String returns the string representation of the value.
func (v *Value) String() string {
	return v.Text()
}
