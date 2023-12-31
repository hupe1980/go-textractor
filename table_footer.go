package textractor

import "strings"

type TableFooter struct {
	base
	words []*Word
}

func (tf *TableFooter) Text() string {
	texts := make([]string, len(tf.words))
	for i, w := range tf.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}
