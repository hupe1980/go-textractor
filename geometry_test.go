package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestBoundingBox(t *testing.T) {
	t.Run("Intersection", func(t *testing.T) {
		t.Run("Bounding boxes do not intersect", func(t *testing.T) {
			bb1 := NewBoundingBox(&types.BoundingBox{Left: 0, Top: 0, Width: 5, Height: 5})
			bb2 := NewBoundingBox(&types.BoundingBox{Left: 10, Top: 10, Width: 5, Height: 5})

			intersection := bb1.Intersection(bb2)

			assert.Nil(t, intersection, "Intersection should be nil when bounding boxes do not intersect")
		})

		t.Run("One bounding box is completely inside the other", func(t *testing.T) {
			bb1 := NewBoundingBox(&types.BoundingBox{Left: 0, Top: 0, Width: 10, Height: 10})
			bb2 := NewBoundingBox(&types.BoundingBox{Left: 3, Top: 3, Width: 3, Height: 3})

			intersection := bb1.Intersection(bb2)

			assert.Equal(t, bb2, intersection, "Intersection should be the smaller bounding box")
		})

		t.Run("Bounding boxes partially intersect", func(t *testing.T) {
			bb1 := NewBoundingBox(&types.BoundingBox{Left: 0, Top: 0, Width: 5, Height: 5})
			bb2 := NewBoundingBox(&types.BoundingBox{Left: 3, Top: 3, Width: 5, Height: 5})

			intersection := bb1.Intersection(bb2)

			assert.NotNil(t, intersection, "Intersection should not be nil")
			assert.InDelta(t, 3, intersection.Left(), 0.001, "Left coordinate of intersection")
			assert.InDelta(t, 3, intersection.Top(), 0.001, "Top coordinate of intersection")
			assert.InDelta(t, 2, intersection.Width(), 0.001, "Width of intersection")
			assert.InDelta(t, 2, intersection.Height(), 0.001, "Height of intersection")
		})
	})

	t.Run("Union", func(t *testing.T) {
		t.Run("Bounding boxes do not intersect", func(t *testing.T) {
			bb1 := NewBoundingBox(&types.BoundingBox{Left: 0, Top: 0, Width: 5, Height: 5})
			bb2 := NewBoundingBox(&types.BoundingBox{Left: 10, Top: 10, Width: 5, Height: 5})

			union := bb1.Union(bb2)

			assert.NotNil(t, union, "Union should not be nil")
			assert.InDelta(t, 0, union.Left(), 0.001, "Left coordinate of union")
			assert.InDelta(t, 0, union.Top(), 0.001, "Top coordinate of union")
			assert.InDelta(t, 15, union.Width(), 0.001, "Width of union")
			assert.InDelta(t, 15, union.Height(), 0.001, "Height of union")
		})

		t.Run("One bounding box is completely inside the other", func(t *testing.T) {
			bb1 := NewBoundingBox(&types.BoundingBox{Left: 0, Top: 0, Width: 10, Height: 10})
			bb2 := NewBoundingBox(&types.BoundingBox{Left: 3, Top: 3, Width: 3, Height: 3})

			union := bb1.Union(bb2)

			assert.Equal(t, bb1, union, "Union should be the larger bounding box")
		})

		t.Run("Bounding boxes partially intersect", func(t *testing.T) {
			bb1 := NewBoundingBox(&types.BoundingBox{Left: 0, Top: 0, Width: 5, Height: 5})
			bb2 := NewBoundingBox(&types.BoundingBox{Left: 3, Top: 3, Width: 5, Height: 5})

			union := bb1.Union(bb2)

			assert.NotNil(t, union, "Union should not be nil")
			assert.InDelta(t, 0, union.Left(), 0.001, "Left coordinate of union")
			assert.InDelta(t, 0, union.Top(), 0.001, "Top coordinate of union")
			assert.InDelta(t, 8, union.Width(), 0.001, "Width of union")
			assert.InDelta(t, 8, union.Height(), 0.001, "Height of union")
		})
	})
}

func TestOrientation(t *testing.T) {
	t.Run("Radians", func(t *testing.T) {
		point0 := NewPoint(types.Point{X: 0, Y: 0})
		point1 := NewPoint(types.Point{X: 1, Y: 1})
		orientation := NewOrientation(point0, point1)
		assert.InDelta(t, 0.785, orientation.Radians(), 0.001, "Radians should be approximately 0.785")
	})

	t.Run("Degrees", func(t *testing.T) {
		point0 := NewPoint(types.Point{X: 0, Y: 0})
		point1 := NewPoint(types.Point{X: 1, Y: 1})
		orientation := NewOrientation(point0, point1)
		assert.InDelta(t, 45.0, orientation.Degrees(), 0.001, "Degrees should be approximately 45.0")
	})
}
