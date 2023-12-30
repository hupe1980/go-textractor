package textractor

import (
	"slices"
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

type Page struct {
	id         string
	number     int
	width      float32
	height     float32
	childIDs   []string
	words      []*Word
	lines      []*Line
	keyValues  []*KeyValue
	tables     []*Table
	layouts    []Layout
	queries    []*Query
	signatures []*Signature
}

func (p *Page) ID() string {
	return p.id
}

func (p *Page) Number() int {
	return p.number
}

func (p *Page) Width() float32 {
	return p.width
}

func (p *Page) Height() float32 {
	return p.height
}

func (p *Page) Words() []*Word {
	return p.words
}

func (p *Page) Lines() []*Line {
	return p.lines
}

func (p *Page) Tables() []*Table {
	return p.tables
}

func (p *Page) KeyValues() []*KeyValue {
	return p.keyValues
}

func (p *Page) Queries() []*Query {
	return p.queries
}

func (p *Page) Signatures() []*Signature {
	return p.signatures
}

func (p *Page) Text(optFns ...func(*TextLinearizationOptions)) string {
	text, _ := p.TextAndWords(optFns...)
	return text
}

func (p *Page) TextAndWords(optFns ...func(*TextLinearizationOptions)) (string, []*Word) {
	pageTexts := make([]string, len(p.layouts))
	wordLists := make([][]*Word, len(p.layouts))

	for i, l := range p.layouts {
		text, words := l.TextAndWords(optFns...)

		pageTexts[i] = text
		wordLists[i] = words
	}

	return strings.Join(pageTexts, "\n"), internal.Concatenate(wordLists...)
}

func (p *Page) SearchValueByKey(key string) []*KeyValue {
	searchKey := strings.ToLower(key)

	var result []*KeyValue

	for _, kv := range p.keyValues {
		if key := kv.Key(); key != nil {
			if strings.Contains(strings.ToLower(key.Text()), searchKey) {
				result = append(result, kv)
			}
		}
	}

	return result
}

func (p *Page) isChild(id string) bool {
	return slices.Contains(p.childIDs, id)
}
