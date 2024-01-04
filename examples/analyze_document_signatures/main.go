package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor"
)

func main() {
	// source: https://aws.amazon.com/blogs/machine-learning/detect-signatures-on-documents-or-images-using-the-signatures-feature-in-amazon-textract/
	file, err := os.Open("examples/analyze_document_signatures/testfile.png")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	cfg, _ := config.LoadDefaultConfig(context.Background())
	client := textract.NewFromConfig(cfg)

	output, err := client.AnalyzeDocument(context.Background(), &textract.AnalyzeDocumentInput{
		Document: &types.Document{
			Bytes: b,
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeSignatures, types.FeatureTypeForms,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	doc, err := textractor.ParseDocumentAPIOutput(&textractor.DocumentAPIOutput{
		DocumentMetadata: output.DocumentMetadata,
		Blocks:           output.Blocks,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over elements in the document
	for _, p := range doc.Pages() {
		for _, s := range p.Signatures() {
			fmt.Printf("ID: %s\n", s.ID())
			fmt.Printf("BoundingBox: %s\n", s.BoundingBox())
			fmt.Printf("Polygon: %s\n", s.Polygon())
			fmt.Println()
		}

		fmt.Println("Search Fields:")

		for _, f := range p.SearchValueByKey("Signature") {
			if k := f.Key(); k != nil {
				fmt.Printf("Key: %s\n", k)
			}

			if v := f.Value(); v != nil {
				fmt.Printf("Value: %s\n", v)
			}

			fmt.Println()
		}
	}
}
