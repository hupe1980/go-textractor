package textractor

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/stretchr/testify/assert"
)

func TestSelectionElement(t *testing.T) {
	t.Run("Status", func(t *testing.T) {
		// Create a SelectionElement with a specific status
		se := &SelectionElement{status: types.SelectionStatusSelected}

		// Check if the Status method returns the correct status
		assert.Equal(t, types.SelectionStatusSelected, se.Status())
	})

	t.Run("IsSelected", func(t *testing.T) {
		// Create a SelectionElement with selected status
		seSelected := &SelectionElement{status: types.SelectionStatusSelected}

		// Create a SelectionElement with a different status
		seNotSelected := &SelectionElement{status: types.SelectionStatusNotSelected}

		// Check if IsSelected returns true for selected status and false otherwise
		assert.True(t, seSelected.IsSelected())
		assert.False(t, seNotSelected.IsSelected())
	})

	t.Run("Text", func(t *testing.T) {
		// Create a SelectionElement with selected status
		seSelected := &SelectionElement{status: types.SelectionStatusSelected}

		// Create a SelectionElement with a different status
		seNotSelected := &SelectionElement{status: types.SelectionStatusNotSelected}

		// Test with default linearization options
		assert.Equal(t, DefaultLinerizationOptions.SelectionElementSelected, seSelected.Text())
		assert.Equal(t, DefaultLinerizationOptions.SelectionElementNotSelected, seNotSelected.Text())

		// Test with custom linearization options
		customOptions := TextLinearizationOptions{
			SelectionElementSelected:    "CUSTOM_SELECTED",
			SelectionElementNotSelected: "CUSTOM_NOT_SELECTED",
		}

		assert.Equal(t, customOptions.SelectionElementSelected, seSelected.Text(func(opts *TextLinearizationOptions) {
			*opts = customOptions
		}))
		assert.Equal(t, customOptions.SelectionElementNotSelected, seNotSelected.Text(func(opts *TextLinearizationOptions) {
			*opts = customOptions
		}))
	})
}
