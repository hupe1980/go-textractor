package textractor

import (
	"math"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// BoundingBox represents the bounding box of a geometry.
type BoundingBox struct {
	boundingBox *types.BoundingBox
}

// NewBoundingBox creates a new BoundingBox instance.
func NewBoundingBox(boundingBox *types.BoundingBox) *BoundingBox {
	return &BoundingBox{
		boundingBox: boundingBox,
	}
}

// Bottom returns the bottom coordinate of the bounding box.
func (bb *BoundingBox) Bottom() float32 {
	return bb.Top() + bb.Height()
}

// HorizontalCenter returns the horizontal center coordinate of the bounding box.
func (bb *BoundingBox) HorizontalCenter() float32 {
	return bb.Left() + bb.Width()/2
}

// Height returns the height of the bounding box.
func (bb *BoundingBox) Height() float32 {
	return bb.boundingBox.Height
}

// Left returns the left coordinate of the bounding box.
func (bb *BoundingBox) Left() float32 {
	return bb.boundingBox.Left
}

// Top returns the top coordinate of the bounding box.
func (bb *BoundingBox) Top() float32 {
	return bb.boundingBox.Top
}

// Right returns the right coordinate of the bounding box.
func (bb *BoundingBox) Right() float32 {
	return bb.Left() + bb.Width()
}

// VerticalCenter returns the vertical center coordinate of the bounding box.
func (bb *BoundingBox) VerticalCenter() float32 {
	return bb.Top() + bb.Height()/2
}

// Width returns the width of the bounding box.
func (bb *BoundingBox) Width() float32 {
	return bb.boundingBox.Width
}

// Union returns a new bounding box that represents the union of two bounding boxes.
func (bb *BoundingBox) Union(other *BoundingBox) *BoundingBox {
	left := float32(math.Min(float64(bb.Left()), float64(other.Left())))
	top := float32(math.Min(float64(bb.Top()), float64(other.Top())))
	right := float32(math.Max(float64(bb.Right()), float64(other.Right())))
	bottom := float32(math.Max(float64(bb.Bottom()), float64(other.Bottom())))

	return NewBoundingBox(&types.BoundingBox{
		Height: bottom - top,
		Left:   left,
		Top:    top,
		Width:  right - left,
	})
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
		return NewBoundingBox(&types.BoundingBox{
			Height: vbottom - vtop,
			Left:   hleft,
			Top:    vtop,
			Width:  hright - hleft,
		})
	}

	return nil
}

// Point represents a 2D point.
type Point struct {
	point types.Point
}

// NewPoint creates a new Point instance.
func NewPoint(point types.Point) *Point {
	return &Point{
		point: point,
	}
}

// X returns the X coordinate of the point.
func (p *Point) X() float32 {
	return p.point.X
}

// Y returns the Y coordinate of the point.
func (p *Point) Y() float32 {
	return p.point.Y
}

// Geometry represents the geometric properties of an element.
type Geometry struct {
	boundingBox *BoundingBox
	polygon     []*Point
}

// NewGeometry creates a new Geometry instance.
func NewGeometry(geometry *types.Geometry) *Geometry {
	polygon := make([]*Point, len(geometry.Polygon))
	for i, p := range geometry.Polygon {
		polygon[i] = NewPoint(p)
	}

	return &Geometry{
		boundingBox: NewBoundingBox(geometry.BoundingBox),
		polygon:     polygon,
	}
}

// BoundingBox returns the bounding box of the geometry.
func (g *Geometry) BoundingBox() *BoundingBox {
	return g.boundingBox
}

// Polygon returns the polygon of the geometry.
func (g *Geometry) Polygon() []*Point {
	return g.polygon
}

// Orientation represents the orientation of a geometric element.
type Orientation struct {
	point0 *Point
	point1 *Point
}

// NewOrientation creates a new Orientation instance.
func NewOrientation(point0, point1 *Point) *Orientation {
	return &Orientation{
		point0: point0,
		point1: point1,
	}
}

// Radians returns the orientation in radians.
func (o *Orientation) Radians() float32 {
	return float32(math.Atan2(float64(o.point1.Y()-o.point0.Y()), float64(o.point1.X()-o.point0.X())))
}

// Degrees returns the orientation in degrees.
func (o *Orientation) Degrees() float32 {
	return (o.Radians() * 180) / math.Pi
}

// Orientation returns the orientation of the geometry.
func (g *Geometry) Orientation() *Orientation {
	if len(g.polygon) < 2 {
		return nil
	}

	return NewOrientation(g.polygon[0], g.polygon[1])
}
