package textractor

import "github.com/aws/aws-sdk-go-v2/service/textract/types"

// Signature represents a signature in a document.
type Signature struct {
	content
}

// NewSignature creates a new Signature instance.
func NewSignature(block types.Block) *Signature {
	return &Signature{
		content: content{block},
	}
}
