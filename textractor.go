// Package textractor provides functionality to work with the Amazon Textract service.
package textractor

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type AnalyzeDocumentPage struct {
	Blocks []types.Block
}

type AnalyzeDocumentSchema struct {
	DocumentMetadata struct {
		Pages *int32 `json:"Pages"`
	} `json:"DocumentMetadata"`
	Blocks []struct {
		BlockType   string   `json:"BlockType"`
		ColumnIndex *int32   `json:"ColumnIndex"`
		ColumnSpan  *int32   `json:"ColumnSpan"`
		ID          *string  `json:"Id"`
		Confidence  *float32 `json:"Confidence"`
		Text        *string  `json:"Text"`
		EntityTypes []string `json:"EntityTypes"`
		Geometry    struct {
			BoundingBox struct {
				Width  float32 `json:"Width"`
				Height float32 `json:"Height"`
				Left   float32 `json:"Left"`
				Top    float32 `json:"Top"`
			} `json:"BoundingBox"`
			Polygon []struct {
				X float32 `json:"X"`
				Y float32 `json:"Y"`
			} `json:"Polygon"`
		} `json:"Geometry"`
		Relationships []struct {
			Type string   `json:"Type"`
			IDs  []string `json:"Ids"`
		} `json:"Relationships"`
	} `json:"Blocks"`
}

func NewAnalyzeDocumentPageFromJSON(data []byte) (*AnalyzeDocumentPage, error) {
	res := new(AnalyzeDocumentSchema)
	if err := json.Unmarshal(data, res); err != nil {
		return nil, err
	}

	rp := new(AnalyzeDocumentPage)

	for _, b := range res.Blocks {
		relationships := []types.Relationship{}
		for _, i := range b.Relationships {
			relationships = append(relationships, types.Relationship{
				Type: types.RelationshipType(i.Type),
				Ids:  i.IDs,
			})
		}

		polygons := []types.Point{}
		for _, p := range b.Geometry.Polygon {
			polygons = append(polygons, types.Point{
				X: p.X,
				Y: p.Y,
			})
		}

		entityTypes := []types.EntityType{}
		for _, et := range b.EntityTypes {
			entityTypes = append(entityTypes, types.EntityType(et))
		}

		rp.Blocks = append(rp.Blocks, types.Block{
			BlockType:     types.BlockType(b.BlockType),
			ColumnIndex:   b.ColumnIndex,
			ColumnSpan:    b.ColumnSpan,
			Id:            b.ID,
			Confidence:    b.Confidence,
			Text:          b.Text,
			EntityTypes:   entityTypes,
			Relationships: relationships,
			Geometry: &types.Geometry{
				BoundingBox: &types.BoundingBox{
					Height: b.Geometry.BoundingBox.Height,
					Width:  b.Geometry.BoundingBox.Width,
					Top:    b.Geometry.BoundingBox.Top,
					Left:   b.Geometry.BoundingBox.Left,
				},
				Polygon: polygons,
			},
		})
	}

	return rp, nil
}

type AnalyzeExpensePage struct {
	ExpenseDocuments []types.ExpenseDocument
}

type AnalyzeExpenseJSONResponse struct {
	DocumentMetadata struct {
		Pages *int32 `json:"Pages"`
	} `json:"DocumentMetadata"`
	ExpenseDocuments []struct {
		LineItemGroups []struct {
			LineItemGroupIndex *int32     `json:"LineItemGroupIndex"`
			LineItems          []struct{} `json:"LineItems"`
		} `json:"LineItemGroups"`
		SummaryFields []struct{} `json:"SummaryFields"`
	} `json:"ExpenseDocuments"`
}
