package textractor

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/textract/types"
)

type LayoutChild interface {
	ID() string
	TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word)
	BoundingBox() *BoundingBox
}

type Layout struct {
	base
	children   []LayoutChild
	noNewLines bool
}

func (l *Layout) AddChildren(children ...LayoutChild) {
	l.children = append(l.children, children...)
}

func (l *Layout) Text(optFns ...func(*TextLinearizationOptions)) string {
	text, _ := l.TextAndWords(optFns...)
	return text
}

func (l *Layout) TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word) {
	opts := DefaultLinerizationOptions

	for _, fn := range optFns {
		fn(&opts)
	}

	if l.BlockType() == types.BlockTypeLayoutHeader && opts.HideHeaderLayout {
		return "", nil
	}

	if l.BlockType() == types.BlockTypeLayoutFooter && opts.HideFooterLayout {
		return "", nil
	}

	if l.BlockType() == types.BlockTypeLayoutFigure && opts.HideFigureLayout {
		return "", nil
	}

	if l.BlockType() == types.BlockTypeLayoutPageNumber && opts.HidePageNumberLayout {
		return "", nil
	}

	text := ""
	words := make([]*Word, 0)
	first := true

	sort.Slice(l.children, func(i, j int) bool {
		return l.children[i].BoundingBox().Top() < l.children[j].BoundingBox().Top()
	})

	for _, group := range groupElementsHorizontally(l.children, opts.HeuristicOverlapRatio) {
		sort.Slice(group, func(i, j int) bool {
			return group[i].BoundingBox().Left() < group[j].BoundingBox().Left()
		})

		for i, child := range group {
			childText, childWords := child.TextAndWords()
			words = append(words, childWords...)

			if l.BlockType() == types.BlockTypeLayoutTable && len(childWords) > 0 {
				columnSep := ""
				if i > 0 {
					columnSep = opts.TableColumnSeparator
				}

				text += columnSep + childText
				// } else if l.BlockType() == types.BlockTypeLayoutKeyValue && len(childWords) > 0 {
				// 	if opts.AddPrefixesAndSuffixesInText {
				// 		text += fmt.Sprintf("%s%s%s", opts.KeyValueLayoutPrefix, text, opts.KeyValueLayoutSuffix)
				// 	}
			} else { // nolint wsl
				sep := ""
				if !first {
					sep = opts.LayoutElementSeparator
				}

				text += sep + childText
			}

			first = false
		}

		if l.BlockType() == types.BlockTypeLayoutTable {
			text += opts.TableRowSeparator
		}
	}

	switch l.BlockType() { // nolint exhaustive
	case types.BlockTypeLayoutPageNumber:
		if opts.AddPrefixesAndSuffixesInText {
			text = fmt.Sprintf("%s%s%s", opts.PageNumberPrefix, text, opts.PageNumberSuffix)
		}
	case types.BlockTypeLayoutTitle:
		if opts.AddPrefixesAndSuffixesInText {
			text = fmt.Sprintf("%s%s%s", opts.TitlePrefix, text, opts.TitleSuffix)
		}
	case types.BlockTypeLayoutSectionHeader:
		if opts.AddPrefixesAndSuffixesInText {
			text = fmt.Sprintf("%s%s%s", opts.SectionHeaderPrefix, text, opts.SectionHeaderSuffix)
		}
	}

	if l.noNewLines {
		// Replace all occurrences of \n with a space
		text = strings.ReplaceAll(text, "\n", " ")

		// Replace consecutive spaces with a single space
		for strings.Contains(text, "  ") {
			text = strings.ReplaceAll(text, "  ", " ")
		}
	}

	return text, words
}

// groupElementsHorizontally groups elements horizontally based on their vertical positions.
// It takes a slice of elements and an overlap ratio as parameters, and returns a 2D slice of grouped elements.
func groupElementsHorizontally(elements []LayoutChild, overlapRatio float32) [][]LayoutChild {
	// Create a copy of the elements to avoid modifying the original slice
	sortedElements := make([]LayoutChild, len(elements))
	copy(sortedElements, elements)

	// Sort elements based on the top position of their bounding boxes
	sort.Slice(sortedElements, func(i, j int) bool {
		return sortedElements[i].BoundingBox().Top() < sortedElements[j].BoundingBox().Top()
	})

	var groupedElements [][]LayoutChild

	// Check if the sorted elements slice is empty
	if len(sortedElements) == 0 {
		return groupedElements
	}

	// verticalOverlap calculates the vertical overlap between two children
	verticalOverlap := func(child1, child2 LayoutChild) float64 {
		t1 := float64(child1.BoundingBox().Top())
		h1 := float64(child1.BoundingBox().Height())
		t2 := float64(child2.BoundingBox().Top())
		h2 := float64(child2.BoundingBox().Height())

		top := math.Max(t1, t2)
		bottom := math.Min(t1+h1, t2+h2)

		return math.Max(bottom-top, 0)
	}

	// shouldGroup determines whether a line should be grouped with an existing group of lines
	shouldGroup := func(child LayoutChild, group []LayoutChild) bool {
		if len(group) == 0 {
			return false
		}

		maxHeight := 0.0
		for _, l := range group {
			maxHeight = math.Max(maxHeight, float64(l.BoundingBox().Height()))
		}

		totalOverlap := 0.0
		for _, l := range group {
			totalOverlap += verticalOverlap(child, l)
		}

		return totalOverlap/maxHeight >= float64(overlapRatio)
	}

	// Initialize the first group with the first element
	currentGroup := []LayoutChild{sortedElements[0]}

	// Iterate through the sorted elements and group them horizontally
	for _, element := range sortedElements[1:] {
		if shouldGroup(element, currentGroup) {
			currentGroup = append(currentGroup, element)
		} else {
			groupedElements = append(groupedElements, currentGroup)
			currentGroup = []LayoutChild{element}
		}
	}

	// Add the last group to the result
	groupedElements = append(groupedElements, currentGroup)

	return groupedElements
}
