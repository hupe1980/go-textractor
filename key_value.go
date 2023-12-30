package textractor

import (
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
	return NewEnclosingBoundingBox(kv.Key().BoundingBox(), kv.Value().BoundingBox())
}

func (kv *KeyValue) Polygon() []*Point {
	// TODO
	panic("not implemented")
}

type Key struct {
	base
	words []*Word
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

func (v *Value) Text() string {
	texts := make([]string, len(v.words))
	for i, w := range v.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

// String returns the string representation of the value.
func (v *Value) String() string {
	return v.Text()
}
