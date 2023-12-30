package textractor

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type DocumentAPIOutput struct {
	DocumentMetadata *types.DocumentMetadata `json:"DocumentMetadata"`
	Blocks           []types.Block           `json:"Blocks"`
}

func ParseDocumentAPIOutput(output *DocumentAPIOutput) (*Document, error) {
	parser := newBlockParser(output.Blocks)

	document := parser.createDocument()

	if len(document.pages) != int(aws.ToInt32(output.DocumentMetadata.Pages)) {
		return nil, fmt.Errorf("number of pages %d does not match metadata %d", len(document.pages), aws.ToInt32(output.DocumentMetadata.Pages))
	}

	return document, nil
}

type AnalyzeIDOutput struct {
	DocumentMetadata  *types.DocumentMetadata  `json:"DocumentMetadata"`
	IdentityDocuments []types.IdentityDocument `json:"IdentityDocuments"`
}

func ParseAnalyzeIDOutput(output *AnalyzeIDOutput) ([]*IdentityDocument, error) {
	parsedIdentityDocuments := make([]*IdentityDocument, len(output.IdentityDocuments))

	for i, d := range output.IdentityDocuments {
		parser := newIdentityDocumentParser(d)
		parsedIdentityDocuments[i] = parser.CreateIdentityDocument()
	}

	if len(parsedIdentityDocuments) != int(aws.ToInt32(output.DocumentMetadata.Pages)) {
		return nil, fmt.Errorf("number of pages %d does not match metadata %d", len(parsedIdentityDocuments), aws.ToInt32(output.DocumentMetadata.Pages))
	}

	return parsedIdentityDocuments, nil
}
