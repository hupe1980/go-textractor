package textractor

import (
	"strings"

	"github.com/hupe1980/go-textractor/internal"
)

type Document struct {
	pages []*Page
}

func (d *Document) Pages() []*Page {
	return d.pages
}

func (d *Document) Words() []*Word {
	words := make([][]*Word, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		words = append(words, p.Words())
	}

	return internal.Concatenate(words...)
}

func (d *Document) Lines() []*Line {
	lines := make([][]*Line, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		lines = append(lines, p.Lines())
	}

	return internal.Concatenate(lines...)
}

func (d *Document) Tables() []*Table {
	tables := make([][]*Table, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		tables = append(tables, p.Tables())
	}

	return internal.Concatenate(tables...)
}

func (d *Document) KeyValues() []*KeyValue {
	keyValues := make([][]*KeyValue, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		keyValues = append(keyValues, p.KeyValues())
	}

	return internal.Concatenate(keyValues...)
}

func (d *Document) Signatures() []*Signature {
	signatures := make([][]*Signature, 0, len(d.Pages()))

	for _, p := range d.Pages() {
		signatures = append(signatures, p.Signatures())
	}

	return internal.Concatenate(signatures...)
}

func (d *Document) Text(optFns ...func(*TextLinearizationOptions)) string {
	pageTexts := make([]string, len(d.Pages()))

	for i, p := range d.Pages() {
		pageTexts[i] = p.Text(optFns...)
	}

	return strings.Join(pageTexts, "\n")
}
