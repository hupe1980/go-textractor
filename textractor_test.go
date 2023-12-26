package textractor

import (
	"io"
	"os"
)

func loadTestdata(filename string) (*AnalyzeDocumentPage, error) { //nolint unparam
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return NewAnalyzeDocumentPageFromJSON(b)
}
