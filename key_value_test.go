package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestKeyValue(t *testing.T) {
	t.Run("Key and Value methods", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}}
		value := &Value{words: []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}}
		kv := &KeyValue{key: key, value: value}

		// Test Key() method
		assert.Equal(t, key, kv.Key())

		// Test Value() method
		assert.Equal(t, value, kv.Value())
	})

	t.Run("Confidence method", func(t *testing.T) {
		// Setup
		key := &Key{base: base{confidence: 0.8}}
		value := &Value{base: base{confidence: 0.7}}
		kv := &KeyValue{key: key, value: value}

		// Test Confidence() method
		expectedConfidence := (0.8 + 0.7) / 2.0
		assert.Equal(t, expectedConfidence, kv.Confidence())
	})

	t.Run("OCRConfidence method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{base: base{confidence: 0.8}}, {base: base{confidence: 0.9}}}}
		value := &Value{words: []*Word{{base: base{confidence: 0.7}}, {base: base{confidence: 0.6}}}}
		kv := &KeyValue{key: key, value: value}

		// Test OCRConfidence() method
		expectedOCRConfidence := &OCRConfidence{
			mean: (0.8 + 0.9 + 0.7 + 0.6) / 4.0,
			min:  0.6,
			max:  0.9,
		}

		assert.Equal(t, expectedOCRConfidence, kv.OCRConfidence())
	})

	t.Run("Words method", func(t *testing.T) {
		// Setup
		keyWords := []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}
		valueWords := []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}
		kv := &KeyValue{key: &Key{words: keyWords}, value: &Value{words: valueWords}}

		// Test Words() method
		expectedWords := append(keyWords, valueWords...)
		assert.Equal(t, expectedWords, kv.Words())
	})

	t.Run("Text method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}}
		value := &Value{words: []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}}
		kv := &KeyValue{key: key, value: value}

		// Test Text() method
		expectedText := "KeyWord1 KeyWord2 ValueWord1 ValueWord2"
		assert.Equal(t, expectedText, kv.Text())
	})

	t.Run("String method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}}
		value := &Value{words: []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}}
		kv := &KeyValue{key: key, value: value}

		// Test String() method
		expectedString := "KeyWord1 KeyWord2 : ValueWord1 ValueWord2"
		assert.Equal(t, expectedString, kv.String())
	})
}

func TestKey(t *testing.T) {
	t.Run("Words method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}}

		// Test Words() method
		assert.Equal(t, key.words, key.Words())
	})

	t.Run("Text method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}}

		// Test Text() method
		expectedText := "KeyWord1 KeyWord2"
		assert.Equal(t, expectedText, key.Text())
	})

	t.Run("OCRConfidence method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{base: base{confidence: 0.8}}, {base: base{confidence: 0.9}}}}

		// Test OCRConfidence() method
		expectedOCRConfidence := &OCRConfidence{
			mean: (0.8 + 0.9) / 2.0,
			min:  0.8,
			max:  0.9,
		}

		actual := key.OCRConfidence()

		assert.InDelta(t, expectedOCRConfidence.Mean(), actual.Mean(), 0.0000001)
		assert.Equal(t, expectedOCRConfidence.Min(), actual.Min())
		assert.Equal(t, expectedOCRConfidence.Max(), actual.Max())
	})

	t.Run("String method", func(t *testing.T) {
		// Setup
		key := &Key{words: []*Word{{text: "KeyWord1"}, {text: "KeyWord2"}}}

		// Test String() method
		expectedString := "KeyWord1 KeyWord2"
		assert.Equal(t, expectedString, key.String())
	})
}

func TestValue(t *testing.T) {
	t.Run("Words and SelectionElement methods", func(t *testing.T) {
		// Setup
		words := []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}
		selectionElement := &SelectionElement{status: types.SelectionStatusSelected}
		value := &Value{words: words, selectionElement: selectionElement}

		// Test Words() method
		assert.Equal(t, words, value.Words())

		// Test SelectionElement() method
		assert.Equal(t, selectionElement, value.SelectionElement())
	})

	t.Run("Text method", func(t *testing.T) {
		// Setup
		words := []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}
		value := &Value{words: words}

		// Test Text() method
		expectedText := "ValueWord1 ValueWord2"
		assert.Equal(t, expectedText, value.Text())
	})

	t.Run("OCRConfidence method", func(t *testing.T) {
		// Setup
		words := []*Word{{base: base{confidence: 0.7}}, {base: base{confidence: 0.6}}}
		selectionElement := &SelectionElement{base: base{confidence: 0.8}}
		value := &Value{words: words, selectionElement: selectionElement}

		// Test OCRConfidence() method
		expectedOCRConfidence := &OCRConfidence{
			mean: 0.8,
			min:  0.8,
			max:  0.8,
		}
		assert.Equal(t, expectedOCRConfidence, value.OCRConfidence())
	})

	t.Run("String method", func(t *testing.T) {
		// Setup
		words := []*Word{{text: "ValueWord1"}, {text: "ValueWord2"}}
		value := &Value{words: words}

		// Test String() method
		expectedString := "ValueWord1 ValueWord2"
		assert.Equal(t, expectedString, value.String())
	})
}
