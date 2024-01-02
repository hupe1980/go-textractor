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
	Text(optFns ...func(*TextLinearizationOptions)) string
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

	var (
		text  string
		words []*Word
		prev  LayoutChild
	)

	for _, group := range groupElementsHorizontally(l.children, opts.HeuristicOverlapRatio) {
		sort.Slice(group, func(i, j int) bool {
			return group[i].BoundingBox().Left() < group[j].BoundingBox().Left()
		})

		for i, child := range group {
			childText := child.Text(func(tlo *TextLinearizationOptions) {
				*tlo = opts
			})

			if l.BlockType() == types.BlockTypeLayoutTable {
				columnSep := ""
				if i > 0 {
					columnSep = opts.TableColumnSeparator
				}

				text += columnSep + childText
			} else if l.BlockType() == types.BlockTypeLayoutKeyValue {
				if opts.AddPrefixesAndSuffixes {
					text += fmt.Sprintf("%s%s%s", opts.KeyValueLayoutPrefix, childText, opts.KeyValueLayoutSuffix)
				}
			} else if partOfSameParagraph(prev, child, opts) {
				text += opts.SameParagraphSeparator + childText
			} else {
				sep := ""
				if prev != nil {
					sep = opts.LayoutElementSeparator
				}

				text += sep + childText
			}

			prev = child
		}

		if l.BlockType() == types.BlockTypeLayoutTable {
			text += opts.TableRowSeparator
		}

		prev = &Line{
			base: base{
				boundingBox: NewEnclosingBoundingBox(group...),
			},
		}
	}

	switch l.BlockType() { // nolint exhaustive
	case types.BlockTypeLayoutPageNumber:
		if opts.AddPrefixesAndSuffixes {
			text = fmt.Sprintf("%s%s%s", opts.PageNumberPrefix, text, opts.PageNumberSuffix)
		}
	case types.BlockTypeLayoutTitle:
		if opts.AddPrefixesAndSuffixes {
			text = fmt.Sprintf("%s%s%s", opts.TitlePrefix, text, opts.TitleSuffix)
		}
	case types.BlockTypeLayoutSectionHeader:
		if opts.AddPrefixesAndSuffixes {
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

func partOfSameParagraph(child1, child2 LayoutChild, options TextLinearizationOptions) bool {
	if child1 != nil && child2 != nil {
		return float32(math.Abs(float64(child1.BoundingBox().Left()-child2.BoundingBox().Left()))) <= options.HeuristicHTolerance*child1.BoundingBox().Width() &&
			float32(math.Abs(float64(child1.BoundingBox().Top()-child2.BoundingBox().Top()))) <= options.HeuristicOverlapRatio*float32(math.Min(float64(child1.BoundingBox().Height()), float64(child2.BoundingBox().Height())))
	}

	return false
}
