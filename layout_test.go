package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupElementsHorizontally(t *testing.T) {
	t.Run("EmptyInput", func(t *testing.T) {
		elements := []LayoutChild{}
		overlapRatio := float32(0.5)
		groups := groupElementsHorizontally(elements, overlapRatio)
		assert.Empty(t, groups, "Expected no groups for empty input")
	})

	t.Run("SingleElement", func(t *testing.T) {
		element := &layoutChildMock{boundingBox: &BoundingBox{top: 10, height: 20}}
		elements := []LayoutChild{element}
		groups := groupElementsHorizontally(elements, 0.5)
		assert.Len(t, groups, 1, "Expected one group for a single element")
		assert.ElementsMatch(t, groups[0], elements, "Expected the group to contain the single element")
	})

	t.Run("MultipleElements", func(t *testing.T) {
		element1 := &layoutChildMock{boundingBox: &BoundingBox{top: 10, height: 20}}
		element2 := &layoutChildMock{boundingBox: &BoundingBox{top: 15, height: 20}}
		element3 := &layoutChildMock{boundingBox: &BoundingBox{top: 30, height: 20}}
		elements := []LayoutChild{element1, element2, element3}
		groups := groupElementsHorizontally(elements, 0.5)
		assert.Len(t, groups, 2, "Expected two groups")
		assert.ElementsMatch(t, groups[0], []LayoutChild{element1, element2}, "Expected the first group to contain element1 and element2")
		assert.ElementsMatch(t, groups[1], []LayoutChild{element3}, "Expected the second group to contain element3")
	})
}

type layoutChildMock struct {
	boundingBox *BoundingBox
}

func (lc *layoutChildMock) BoundingBox() *BoundingBox {
	return lc.boundingBox
}

func (lc *layoutChildMock) Words() []*Word {
	return nil
}

func (lc *layoutChildMock) Text(_ ...func(*TextLinearizationOptions)) string {
	return ""
}
