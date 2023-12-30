package textractor

import "github.com/aws/aws-sdk-go-v2/service/textract/types"

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
