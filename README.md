# 📄 go-textractor
![Build Status](https://github.com/hupe1980/go-textractor/workflows/Build/badge.svg) 
[![Go Reference](https://pkg.go.dev/badge/github.com/hupe1980/go-textractor.svg)](https://pkg.go.dev/github.com/hupe1980/go-textractor)
[![goreportcard](https://goreportcard.com/badge/github.com/hupe1980/go-textractor)](https://goreportcard.com/report/github.com/hupe1980/go-textractor)
[![codecov](https://codecov.io/gh/hupe1980/go-textractor/branch/main/graph/badge.svg?token=VEDVMNI1TV)](https://codecov.io/gh/hupe1980/go-textractor)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
> Amazon textract response parser written in go.

## Installation
Use Go modules to include go-textractor in your project:
```
go get github.com/hupe1980/go-textractor
```

## Usage
```golang
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
	file, err := os.Open("example/testfile.pdf")
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

	res, err := client.AnalyzeDocument(context.Background(), &textract.AnalyzeDocumentInput{
		Document: &types.Document{
			Bytes: b,
		},
		FeatureTypes: []types.FeatureType{
			types.FeatureTypeTables,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	doc := textractor.NewDocument(&textractor.ResponsePage{Blocks: res.Blocks})

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

		// Print fields
		for _, f := range p.Form().Fields() {
			fmt.Printf("Field: Key: %s, Value: %s\n", f.Key(), f.Value())
		}
	}
}
```

For more example usage, see [examples](./examples).

## Contributing
Contributions are welcome! Feel free to open an issue or submit a pull request for any improvements or new features you would like to see.

## License
This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.