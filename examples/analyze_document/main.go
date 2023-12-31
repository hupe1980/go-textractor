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
	file, err := os.Open("examples/analyze_document/testfile.pdf")
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
			types.FeatureTypeTables, types.FeatureTypeForms,
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
		// Print lines and words
		for _, l := range p.Lines() {
			fmt.Printf("Line: %s (%f)\n", l.Text(), l.Confidence())

			for _, w := range l.Words() {
				fmt.Printf("Word: %s (%f)\n", w.Text(), w.Confidence())
			}
		}

		// Print tables
		for _, t := range p.Tables() {
			for r, row := range t.Rows() {
				for c, cell := range row.Cells() {
					fmt.Printf("Table[%d][%d] = %s (%f)\n", r, c, cell.Text(), cell.Confidence())
				}
			}
		}

		// Print key values
		for _, kv := range p.KeyValues() {
			fmt.Printf("Key: %s, Value: %s\n", kv.Key(), kv.Value())
		}
	}
}
