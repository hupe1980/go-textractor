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
	// source: https://aws.amazon.com/blogs/machine-learning/announcing-support-for-extracting-data-from-identity-documents-using-amazon-textract/
	file, err := os.Open("examples/analyze_id/testfile.jpeg")
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

	output, err := client.AnalyzeID(context.Background(), &textract.AnalyzeIDInput{
		DocumentPages: []types.Document{
			{
				Bytes: b,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	parser := textractor.NewAnalyzeIDOutputParser(output)

	for _, doc := range parser.Documents() {
		fmt.Printf("Document type: %s\n\n", doc.Type())

		for _, f := range doc.Fields() {
			value := f.Value()

			if f.IsNormalized() {
				date, err := f.NormalizedValue().DateValue()
				if err != nil {
					log.Fatal(err)
				}

				value = date.Format("2006-01-02")
			}

			fmt.Printf("%s: %s\n", f.Type(), value)
		}

		// for _, p := range doc.Pages() {
		// 	fmt.Println(p.Text())
		// }
	}
}
