package textractor

import (
	"fmt"
	"math"
)

type BoundingBox struct {
	height float32
	left   float32
	top    float32
	width  float32
}

// Bottom returns the bottom coordinate of the bounding box.
func (bb *BoundingBox) Bottom() float32 {
	return bb.Top() + bb.Height()
}

// HorizontalCenter returns the horizontal center coordinate of the bounding box.
func (bb *BoundingBox) HorizontalCenter() float32 {
	return bb.Left() + bb.Width()/2
}

func (bb *BoundingBox) Height() float32 {
	return bb.height
}

func (bb *BoundingBox) Left() float32 {
	return bb.left
}

func (bb *BoundingBox) Top() float32 {
	return bb.top
}

func (bb *BoundingBox) Width() float32 {
	return bb.width
}

// Right returns the right coordinate of the bounding box.
func (bb *BoundingBox) Right() float32 {
	return bb.Left() + bb.Width()
}

// VerticalCenter returns the vertical center coordinate of the bounding box.
func (bb *BoundingBox) VerticalCenter() float32 {
	return bb.Top() + bb.Height()/2
}

// Area calculates and returns the area of the bounding box.
// If either the width or height of the bounding box is less than zero,
// the area is considered zero to prevent negative area values.
func (bb *BoundingBox) Area() float32 {
	if bb.Width() < 0 || bb.Height() < 0 {
		return 0
	}

	return bb.Width() * bb.Height()
}

// Intersection returns a new bounding box that represents the intersection of two bounding boxes.
func (bb *BoundingBox) Intersection(other *BoundingBox) *BoundingBox {
	vtop := float32(math.Max(float64(bb.Top()), float64(other.Top())))
	vbottom := float32(math.Min(float64(bb.Bottom()), float64(other.Bottom())))
	visect := float32(math.Max(0, float64(vbottom-vtop)))
	hleft := float32(math.Max(float64(bb.Left()), float64(other.Left())))
	hright := float32(math.Min(float64(bb.Right()), float64(other.Right())))
	hisect := float32(math.Max(0, float64(hright-hleft)))

	if hisect > 0 && visect > 0 {
		return &BoundingBox{
			height: vbottom - vtop,
			left:   hleft,
			top:    vtop,
			width:  hright - hleft,
		}
	}

	return nil
}

// String returns a string representation of the bounding box.
func (bb *BoundingBox) String() string {
	return fmt.Sprintf("width: %f, height: %f, left: %f, top: %f", bb.Width(), bb.Height(), bb.Left(), bb.Top())
}

type BoundingBoxAccessor interface {
	BoundingBox() *BoundingBox
}

// NewEnclosingBoundingBox returns a new bounding box that represents the union of multiple bounding boxes.
func NewEnclosingBoundingBox[T BoundingBoxAccessor](accessors ...T) *BoundingBox {
	if len(accessors) == 0 {
		return nil
	}

	bboxes := make([]*BoundingBox, 0, len(accessors))
	for _, a := range accessors {
		bboxes = append(bboxes, a.BoundingBox())
	}

	left, top, right, bottom := float32(math.Inf(1)), float32(math.Inf(1)), float32(math.Inf(-1)), float32(math.Inf(-1))

	for _, bb := range bboxes {
		if bb == nil {
			continue
		}

		left = float32(math.Min(float64(left), float64(bb.Left())))
		top = float32(math.Min(float64(top), float64(bb.Top())))
		right = float32(math.Max(float64(right), float64(bb.Right())))
		bottom = float32(math.Max(float64(bottom), float64(bb.Bottom())))
	}

	return &BoundingBox{
		height: bottom - top,
		left:   left,
		top:    top,
		width:  right - left,
	}
}

// Point represents a 2D point.
type Point struct {
	x, y float32
}

// X returns the X coordinate of the point.
func (p *Point) X() float32 {
	return p.x
}

// Y returns the Y coordinate of the point.
func (p *Point) Y() float32 {
	return p.y
}

// String returns a string representation of the Point, including its X and Y coordinates.
func (p *Point) String() string {
	return fmt.Sprintf("x: %f, y: %f", p.x, p.y)
}

// Orientation represents the orientation of a geometric element.
type Orientation struct {
	point0 *Point
	point1 *Point
}

// Radians returns the orientation in radians.
func (o *Orientation) Radians() float32 {
	return float32(math.Atan2(float64(o.point1.Y()-o.point0.Y()), float64(o.point1.X()-o.point0.X())))
}

// Degrees returns the orientation in degrees.
func (o *Orientation) Degrees() float32 {
	return (o.Radians() * 180) / math.Pi
}
