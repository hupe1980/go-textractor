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
	file, err := os.Open("examples/analyze_document_for_llm/testfile.png")
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
			types.FeatureTypeLayout, types.FeatureTypeSignatures, types.FeatureTypeForms,
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

	fmt.Println(doc.Text(func(tlo *textractor.TextLinearizationOptions) {
		tlo.SelectionElementSelected = "[X]"
		tlo.SelectionElementNotSelected = "[ ]"
		tlo.SignatureToken = "[SIGNATURE]"
	}))

}
