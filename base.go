package textractor

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// base represents the base information shared among different types of blocks.
type base struct {
	id          string          // Identifier for the block
	confidence  float32         // Confidence for the block
	blockType   types.BlockType // Type of the block
	boundingBox *BoundingBox    // Bounding box information
	polygon     []*Point        // Polygon information
	page        *Page           // Page information
	raw         types.Block     // Raw block data
}

// newBase creates a new base instance from the provided Textract block and page information.
func newBase(b types.Block, p *Page) base {
	polygon := make([]*Point, len(b.Geometry.Polygon))
	for i, p := range b.Geometry.Polygon {
		polygon[i] = &Point{
			x: p.X,
			y: p.Y,
		}
	}

	return base{
		id:         aws.ToString(b.Id),
		confidence: aws.ToFloat32(b.Confidence),
		blockType:  b.BlockType,
		boundingBox: &BoundingBox{
			height: b.Geometry.BoundingBox.Height,
			left:   b.Geometry.BoundingBox.Left,
			top:    b.Geometry.BoundingBox.Top,
			width:  b.Geometry.BoundingBox.Width,
		},
		polygon: polygon,
		page:    p,
		raw:     b,
	}
}

// ID returns the identifier of the block.
func (b *base) ID() string {
	return b.id
}

// Confidence returns the confidence of the block.
func (b *base) Confidence() float32 {
	return b.confidence
}

// BlockType returns the type of the block.
func (b *base) BlockType() types.BlockType {
	return b.blockType
}

// BoundingBox returns the bounding box information of the block.
func (b *base) BoundingBox() *BoundingBox {
	return b.boundingBox
}

// Polygon returns the polygon information of the block.
func (b *base) Polygon() []*Point {
	return b.polygon
}

// PageNumber returns the page number associated with the block.
func (b *base) PageNumber() int {
	return b.page.Number()
}

// Raw returns the raw block data.
func (b *base) Raw() types.Block {
	return b.raw
}
