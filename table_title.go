package textractor

import "strings"

type TableTitle struct {
	base
	words []*Word
}

func (tt *TableTitle) Text() string {
	texts := make([]string, len(tt.words))
	for i, w := range tt.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}
