package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/textract"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor"
)

func main() {
	// source: https://aws.amazon.com/blogs/machine-learning/specify-and-extract-information-from-documents-using-the-new-queries-feature-in-amazon-textract/
	file, err := os.Open("examples/analyze_document_queries/testfile.jpeg")
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
			types.FeatureTypeQueries,
		},
		QueriesConfig: &types.QueriesConfig{
			Queries: []types.Query{
				{
					Text:  aws.String("What is the year to date gross pay?"),
					Alias: aws.String("PAYSTUB_YTD_GROSS"),
				},
				{
					Text:  aws.String("What is the current gross pay?"),
					Alias: aws.String("PAYSTUB_CURRENT_GROSS"),
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	doc := textractor.NewDocument(output.Blocks)

	// Iterate over elements in the document
	for _, p := range doc.Pages() {
		for _, q := range p.Queries() {
			fmt.Printf("Question: %s\n", q.Text())
			fmt.Printf("Alias: %s\n", q.Alias())

			if r := q.TopResult(); r != nil {
				fmt.Printf("Answer: %s\n", r.Text())
			}

			fmt.Println()
		}
	}
}
