package textractor

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

// Queries represents a slice of Query instances.
type Queries []*Query

// QueryResult represents the result of a Textract query.
type QueryResult struct {
	content
}

// NewQueryResult creates a new QueryResult instance.
func NewQueryResult(block types.Block) *QueryResult {
	return &QueryResult{
		content: content{block},
	}
}

// Text returns the text content of the query result.
func (qr *QueryResult) Text() string {
	return aws.ToString(qr.Block().Text)
}

// Query represents a Textract query.
type Query struct {
	block   types.Block
	results []*QueryResult
}

// NewQuery creates a new Query instance.
func NewQuery(block types.Block, blockMap map[string]types.Block) *Query {
	query := &Query{
		block: block,
	}

	for _, r := range block.Relationships {
		if r.Type == types.RelationshipTypeAnswer {
			for _, i := range r.Ids {
				b := blockMap[i]
				if b.BlockType == types.BlockTypeQueryResult {
					query.results = append(query.results, NewQueryResult(b))
				}
			}
		}
	}

	return query
}

// Alias returns the alias of the query.
func (q *Query) Alias() string {
	return aws.ToString(q.block.Query.Alias)
}

// Text returns the text content of the query.
func (q *Query) Text() string {
	return aws.ToString(q.block.Query.Text)
}

// TopResult retrieves the top result by confidence score, if any are available.
func (q *Query) TopResult() *QueryResult {
	r := q.ResultsByConfidence()
	if len(r) > 0 {
		return r[0]
	}

	return nil
}

// ResultsByConfidence lists this query instance's results, sorted from most to least confident.
func (q *Query) ResultsByConfidence() []*QueryResult {
	sortedResults := make([]*QueryResult, len(q.results))
	copy(sortedResults, q.results)
	sort.Slice(sortedResults, func(i, j int) bool {
		// Negative -> a sorted before b
		return sortedResults[j].Confidence() < sortedResults[i].Confidence()
	})

	return sortedResults
}
