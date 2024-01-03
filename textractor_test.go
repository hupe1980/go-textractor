package textractor

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDocumentAPIOutput(t *testing.T) {
	t.Run("Layout table without table", func(t *testing.T) {
		res, err := loadDocumentAPIOutputTestdata("testdata/test-layout-table-without-table.json")
		assert.NoError(t, err)

		doc, err := ParseDocumentAPIOutput(res)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(doc.Tables()))

		text := doc.Text(func(tlo *TextLinearizationOptions) {
			tlo.HideFigureLayout = true
			tlo.TitlePrefix = "# "
			tlo.SectionHeaderPrefix = "## "
		})

		//fmt.Println(text)

		assert.Equal(t, `

(a) Original
(b) Reconstructed
Figure 3. Example for the Learn To Reconstruct task output on the IIT-CDIP dataset
Table 1. Entity-level F1 scores of two entity extraction tasks: FUNSD and CORD.
## 4.4. Ablation Study
We conduct an extensive ablation study using the CORD dataset.
Model	#param (M)	FUNSD	CORD
LayoutLMvl-base	160	79.27	-
LayoutLMvl-large	390	77.89	94.93
LayoutLMv2-base	200	82.76	94.95
TILT-base	230	-	95.11
LayoutLMv2-large	426	84.20	96.01
TILT-large	780	-	96.33
DocFormer-base	183	83.34	96.33
DocFormer-large	533	84.55	96.99
MATrIX (ours)	166	78.60	96.05

`, text[:575])
	})

	t.Run("SimpleTableLayout", func(t *testing.T) {
		res, err := loadDocumentAPIOutputTestdata("testdata/test-simple-table-layout.json")
		//res, err := loadDocumentAPIOutputTestdata("testdata/table-example-response.json")
		assert.NoError(t, err)

		doc, err := ParseDocumentAPIOutput(res)
		assert.NoError(t, err)

		text := doc.Text(func(tlo *TextLinearizationOptions) {
			tlo.TitlePrefix = "# "
			tlo.SectionHeaderPrefix = "## "
			tlo.TableLinearizationFormat = "markdown"
		})

		//fmt.Println(text)

		assert.Equal(t, `# New Document
## Paragraph 1
Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.


| A  |  B  | C  |
|----|-----|----|
| A1 | b1  | C1 |
| A2 | B2  | C2 |
| A3 | BC3 |    |
| A4 | B4  | C4 |

`, text)
	})

	t.Run("ForLLM", func(t *testing.T) {
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
	})

	t.Run("No layout", func(t *testing.T) {
		res, err := loadDocumentAPIOutputTestdata("testdata/test-document.json")
		assert.NoError(t, err)

		doc, err := ParseDocumentAPIOutput(res)
		assert.NoError(t, err)

		//fmt.Println(doc.Text())

		assert.Equal(t, 51, len(doc.Words()))
		assert.Equal(t, 24, len(doc.Lines()))
		assert.Equal(t, 5, len(doc.KeyValues()))
		assert.Equal(t, 1, len(doc.Tables()))
	})
}

func TestParseAnalyzeIDOutput(t *testing.T) {
	res, err := loadAnalyzeIDOutputTestdata("testdata/test-analyze-id-response.json")
	assert.NoError(t, err)

	idocuments, err := ParseAnalyzeIDOutput(res)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(idocuments))
	assert.Equal(t, 21, len(idocuments[0].Fields()))
	assert.Equal(t, IdentityDocumentTypeDriverLicenseFront, idocuments[0].IdentityDocumentType())
	assert.Equal(t, "GARCIA", idocuments[0].FieldByType(IdentityDocumentFieldTypeFirstName).Value())
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

func loadAnalyzeIDOutputTestdata(filename string) (*AnalyzeIDOutput, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	output := new(AnalyzeIDOutput)
	if err := json.Unmarshal(data, output); err != nil {
		return nil, err
	}

	return output, nil
}
