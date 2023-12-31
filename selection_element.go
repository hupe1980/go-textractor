package textractor

import (
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/google/uuid"
)

type SelectionElement struct {
	base
	status types.SelectionStatus
}

// Status returns the selection status of the element.
func (se *SelectionElement) Status() types.SelectionStatus {
	return se.status
}

// IsSelected checks if the element is selected.
func (se *SelectionElement) IsSelected() bool {
	return se.Status() == types.SelectionStatusSelected
}

func (se *SelectionElement) Words() []*Word {
	_, words := se.TextAndWords()
	return words
}

func (se *SelectionElement) TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word) {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	text := opts.SelectionElementNotSelected
	if se.IsSelected() {
		text = opts.SelectionElementSelected
	}

	w := &Word{
		base: base{
			id:          uuid.New().String(),
			confidence:  se.Confidence(),
			blockType:   types.BlockTypeWord,
			boundingBox: se.BoundingBox(),
			page:        se.page,
		},
		text: text,
	}

	w.line = &Line{
		base: base{
			id:          uuid.New().String(),
			confidence:  se.Confidence(),
			blockType:   types.BlockTypeLine,
			boundingBox: se.BoundingBox(),
			page:        se.page,
		},
		words: []*Word{w},
	}

	return text, []*Word{w}
}
