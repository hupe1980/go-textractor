package textractor

import (
	"strings"
)

type Line struct {
	base
	words []*Word
}

func (l *Line) Text() string {
	texts := make([]string, len(l.words))
	for i, w := range l.words {
		texts[i] = w.Text()
	}

	return strings.Join(texts, " ")
}

func (l *Line) Words() []*Word {
	return l.words
}
