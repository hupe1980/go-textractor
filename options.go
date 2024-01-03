package textractor

// TextLinearizationOptions defines how a document is linearized into a text string.
type TextLinearizationOptions struct {
	// LinearizeTables includes tables in the linearized output.
	LinearizeTables bool

	// LinearizeKeyValues includes form key and values in the linearized output.
	LinearizeKeyValues bool

	// RemoveNewLinesInListElements removes new lines in list elements.
	RemoveNewLinesInListElements bool

	// MaxNumberOfConsecutiveNewLines sets the maximum number of consecutive new lines to keep, removing extra whitespace.
	MaxNumberOfConsecutiveNewLines int

	// HideHeaderLayout hides headers in the linearized output.
	HideHeaderLayout bool

	// HideFooterLayout hides footers in the linearized output.
	HideFooterLayout bool

	// HideFigureLayout hides figures in the linearized output.
	HideFigureLayout bool

	// HidePageNumberLayout hides page numbers in the linearized output.
	HidePageNumberLayout bool

	// PageNumberPrefix is the prefix for page number layout elements.
	PageNumberPrefix string

	// PageNumberSuffix is the suffix for page number layout elements.
	PageNumberSuffix string

	// SameParagraphSeparator is the separator to use when combining elements within a text block.
	SameParagraphSeparator string

	// LayoutElementSeparator is the separator to use when combining linearized layout elements.
	LayoutElementSeparator string

	// ListElementSeparator is the separator for elements in a list layout.
	ListElementSeparator string

	// ListLayoutPrefix is the prefix for list layout elements (parent).
	ListLayoutPrefix string

	// ListLayoutSuffix is the suffix for list layout elements (parent).
	ListLayoutSuffix string

	// ListElementPrefix is the prefix for elements in a list layout (children).
	ListElementPrefix string

	// ListElementSuffix is the suffix for elements in a list layout (children).
	ListElementSuffix string

	// TitlePrefix is the prefix for title layout elements.
	TitlePrefix string

	// TitleSuffix is the suffix for title layout elements.
	TitleSuffix string

	// TableLayoutPrefix is the prefix for table elements.
	TableLayoutPrefix string

	// TableLayoutSuffix is the suffix for table elements.
	TableLayoutSuffix string

	// TableLinearizationFormat sets how to represent tables in the linearized output. Choices are plaintext or markdown.
	TableLinearizationFormat string

	// TableMinTableWords is the threshold below which tables will be rendered as words instead of using table layout.
	TableMinTableWords int

	// TableColumnSeparator is the table column separator, used when linearizing layout tables, not used if AnalyzeDocument was called with the TABLES feature.
	TableColumnSeparator string

	// TablePrefix is the prefix for table layout.
	TablePrefix string

	// TableSuffix is the suffix for table layout.
	TableSuffix string

	// TableRowSeparator is the table row separator.
	TableRowSeparator string

	// TableRowPrefix is the prefix for table row.
	TableRowPrefix string

	// TableRowSuffix is the suffix for table row.
	TableRowSuffix string

	// TableCellPrefix is the prefix for table cell.
	TableCellPrefix string

	// TableCellSuffix is the suffix for table cell.
	TableCellSuffix string

	// SectionHeaderPrefix is the prefix for section header layout elements.
	SectionHeaderPrefix string

	// SectionHeaderSuffix is the suffix for section header layout elements.
	SectionHeaderSuffix string

	// TextPrefix is the prefix for text layout elements.
	TextPrefix string

	// TextSuffix is the suffix for text layout elements.
	TextSuffix string

	// KeyValueLayoutPrefix is the prefix for key_value layout elements (not for individual key-value elements).
	KeyValueLayoutPrefix string

	// KeyValueLayoutSuffix is the suffix for key_value layout elements (not for individual key-value elements).
	KeyValueLayoutSuffix string

	// KeyValuePrefix is the prefix for key-value elements.
	KeyValuePrefix string

	// KeyValueSuffix is the suffix for key-value elements.
	KeyValueSuffix string

	// KeyPrefix is the prefix for key elements.
	KeyPrefix string

	// KeySuffix is the suffix for key elements.
	KeySuffix string

	// ValuePrefix is the prefix for value elements.
	ValuePrefix string

	// ValueSuffix is the suffix for value elements.
	ValueSuffix string

	// SelectionElementSelected is the representation for selection elements when selected.
	SelectionElementSelected string

	// SelectionElementNotSelected is the representation for selection elements when not selected.
	SelectionElementNotSelected string

	// HeuristicHTolerance sets how much the line below and above the current line should differ in width to be separated.
	HeuristicHTolerance float64

	// HeuristicLineBreakThreshold sets how much space is acceptable between two lines before splitting them. Expressed in multiple of min heights.
	HeuristicLineBreakThreshold float64

	// HeuristicOverlapRatio sets how much vertical overlap is tolerated between two subsequent lines before merging them into a single line.
	HeuristicOverlapRatio float64

	// SignatureToken is the signature representation in the linearized text.
	SignatureToken string
}

var DefaultLinerizationOptions = TextLinearizationOptions{
	LinearizeTables:                true,
	LinearizeKeyValues:             true,
	RemoveNewLinesInListElements:   true,
	MaxNumberOfConsecutiveNewLines: 2,
	HideHeaderLayout:               false,
	HideFooterLayout:               false,
	HideFigureLayout:               false,
	HidePageNumberLayout:           false,
	PageNumberPrefix:               "",
	PageNumberSuffix:               "",
	SameParagraphSeparator:         " ",
	LayoutElementSeparator:         "\n\n",
	ListElementSeparator:           "\n",
	ListLayoutPrefix:               "",
	ListLayoutSuffix:               "",
	ListElementPrefix:              "",
	ListElementSuffix:              "",
	TitlePrefix:                    "",
	TitleSuffix:                    "",
	TableLayoutPrefix:              "\n\n",
	TableLayoutSuffix:              "\n",
	TableLinearizationFormat:       "plaintext",
	TableMinTableWords:             0,
	TableColumnSeparator:           "\t",
	TablePrefix:                    "",
	TableSuffix:                    "",
	TableRowSeparator:              "\n",
	TableRowPrefix:                 "",
	TableRowSuffix:                 "",
	TableCellPrefix:                "",
	TableCellSuffix:                "",
	SectionHeaderPrefix:            "",
	SectionHeaderSuffix:            "",
	TextPrefix:                     "",
	TextSuffix:                     "",
	KeyValueLayoutPrefix:           "\n\n",
	KeyValueLayoutSuffix:           "",
	KeyValuePrefix:                 "",
	KeyValueSuffix:                 "",
	KeyPrefix:                      "",
	KeySuffix:                      "",
	ValuePrefix:                    "",
	ValueSuffix:                    "",
	SelectionElementSelected:       "[X]",
	SelectionElementNotSelected:    "[ ]",
	HeuristicHTolerance:            0.3,
	HeuristicLineBreakThreshold:    0.9,
	HeuristicOverlapRatio:          0.5,
	SignatureToken:                 "[SIGNATURE]",
}
