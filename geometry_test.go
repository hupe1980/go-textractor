package textractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoundingBox(t *testing.T) {
	t.Run("Area", func(t *testing.T) {
		bb1 := &BoundingBox{left: 0, top: 0, width: 4, height: 3}
		assert.Equal(t, float32(12), bb1.Area(), "Area should be 12 for a bounding box of width 4 and height 3")

		bb2 := &BoundingBox{left: 2, top: 2, width: 0, height: 3}
		assert.Equal(t, float32(0), bb2.Area(), "Area should be 0 for a bounding box with zero width")

		bb3 := &BoundingBox{left: 2, top: 2, width: 4, height: 0}
		assert.Equal(t, float32(0), bb3.Area(), "Area should be 0 for a bounding box with zero height")

		bb4 := &BoundingBox{left: 0, top: 0, width: -2, height: 3}
		assert.Equal(t, float32(0), bb4.Area(), "Area should be 0 for a bounding box with negative width")

		bb5 := &BoundingBox{left: 0, top: 0, width: 2, height: -3}
		assert.Equal(t, float32(0), bb5.Area(), "Area should be 0 for a bounding box with negative height")
	})

	t.Run("Intersection", func(t *testing.T) {
		t.Run("Bounding boxes do not intersect", func(t *testing.T) {
			bb1 := &BoundingBox{left: 0, top: 0, width: 5, height: 5}
			bb2 := &BoundingBox{left: 10, top: 10, width: 5, height: 5}

			intersection := bb1.Intersection(bb2)

			assert.Nil(t, intersection, "Intersection should be nil when bounding boxes do not intersect")
		})

		t.Run("One bounding box is completely inside the other", func(t *testing.T) {
			bb1 := &BoundingBox{left: 0, top: 0, width: 10, height: 10}
			bb2 := &BoundingBox{left: 3, top: 3, width: 3, height: 3}

			intersection := bb1.Intersection(bb2)

			assert.Equal(t, bb2, intersection, "Intersection should be the smaller bounding box")
		})

		t.Run("Bounding boxes partially intersect", func(t *testing.T) {
			bb1 := &BoundingBox{left: 0, top: 0, width: 5, height: 5}
			bb2 := &BoundingBox{left: 3, top: 3, width: 5, height: 5}

			intersection := bb1.Intersection(bb2)

			assert.NotNil(t, intersection, "Intersection should not be nil")
			assert.InDelta(t, 3, intersection.Left(), 0.001, "Left coordinate of intersection")
			assert.InDelta(t, 3, intersection.Top(), 0.001, "Top coordinate of intersection")
			assert.InDelta(t, 2, intersection.Width(), 0.001, "Width of intersection")
			assert.InDelta(t, 2, intersection.Height(), 0.001, "Height of intersection")
		})
	})
}

func TestNewEnclosingBoundingBox(t *testing.T) {
	// Test case 1: Bounding boxes with positive coordinates
	bbox1 := &BoundingBox{left: 0, top: 0, width: 2, height: 2}
	bbox2 := &BoundingBox{left: 1, top: 1, width: 2, height: 2}
	result1 := NewEnclosingBoundingBox(bbox1, bbox2)
	assert.Equal(t, &BoundingBox{left: 0, top: 0, width: 3, height: 3}, result1)

	// Test case 2: Bounding boxes with negative coordinates
	bbox3 := &BoundingBox{left: -2, top: -2, width: 3, height: 3}
	bbox4 := &BoundingBox{left: -3, top: -3, width: 2, height: 2}
	result2 := NewEnclosingBoundingBox(bbox3, bbox4)

	assert.Equal(t, &BoundingBox{left: -3, top: -3, width: 4, height: 4}, result2)

	// Test case 3: Bounding boxes with one nil
	result3 := NewEnclosingBoundingBox(nil, bbox1, bbox2)
	assert.Equal(t, &BoundingBox{left: 0, top: 0, width: 3, height: 3}, result3)

	// Test case 4: Empty input
	result4 := NewEnclosingBoundingBox()
	assert.Nil(t, result4)

	// Test case 5: Bounding boxes with floating-point coordinates
	bbox5 := &BoundingBox{left: 0.1, top: 0.2, width: 2.5, height: 2.8}
	bbox6 := &BoundingBox{left: 1.3, top: 1.5, width: 2.2, height: 2.6}
	result5 := NewEnclosingBoundingBox(bbox5, bbox6)
	expectedResult5 := &BoundingBox{left: 0.1, top: 0.2, width: 3.4, height: 3.9}
	assert.InDelta(t, expectedResult5.Left(), result5.Left(), 0.0001, "Floating-point coordinates not within tolerance")
	assert.InDelta(t, expectedResult5.Top(), result5.Top(), 0.0001, "Floating-point coordinates not within tolerance")
	assert.InDelta(t, expectedResult5.Width(), result5.Width(), 0.0001, "Floating-point coordinates not within tolerance")
	assert.InDelta(t, expectedResult5.Height(), result5.Height(), 0.0001, "Floating-point coordinates not within tolerance")
}

func TestOrientation(t *testing.T) {
	t.Run("Radians", func(t *testing.T) {
		point0 := &Point{x: 0, y: 0}
		point1 := &Point{x: 1, y: 1}
		orientation := &Orientation{point0, point1}
		assert.InDelta(t, 0.785, orientation.Radians(), 0.001, "Radians should be approximately 0.785")
	})

	t.Run("Degrees", func(t *testing.T) {
		point0 := &Point{x: 0, y: 0}
		point1 := &Point{x: 1, y: 1}
		orientation := &Orientation{point0, point1}
		assert.InDelta(t, 45.0, orientation.Degrees(), 0.001, "Degrees should be approximately 45.0")
	})
}
