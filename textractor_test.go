package textractor

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocumentAPIOutputForLLM(t *testing.T) {
	res, err := loadDocumentAPIOutputTestdata("testdata/test-response-for-llm.json")
	assert.NoError(t, err)

	doc, err := ParseDocumentAPIOutput(res)
	assert.NoError(t, err)

	text := doc.Text(func(tlo *TextLinearizationOptions) {
		tlo.SelectionElementSelected = "[X]"
		tlo.SelectionElementNotSelected = "[ ]"
		tlo.SignatureToken = "[SIGNATURE]"
	})
	//fmt.Println(doc.Text())

	sigCount := strings.Count(text, "[SIGNATURE]")
	assert.Equal(t, 3, sigCount)
	assert.Equal(t, 3, len(doc.Signatures()))
}

func TestParseDocumentAPIOutput(t *testing.T) {
	res, err := loadDocumentAPIOutputTestdata("testdata/test-document.json")
	assert.NoError(t, err)

	doc, err := ParseDocumentAPIOutput(res)
	assert.NoError(t, err)

	assert.Equal(t, 51, len(doc.Words()))
	assert.Equal(t, 24, len(doc.Lines()))
	assert.Equal(t, 5, len(doc.KeyValues()))
	assert.Equal(t, 1, len(doc.Tables()))
}

func loadDocumentAPIOutputTestdata(filename string) (*DocumentAPIOutput, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	output := new(DocumentAPIOutput)
	if err := json.Unmarshal(data, output); err != nil {
		return nil, err
	}

	return output, nil
}
