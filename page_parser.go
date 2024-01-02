package textractor

import (
	"fmt"
	"slices"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/google/uuid"
)

const threshold = 0.8

type pageParser struct {
	bp        *blockParser
	page      *Page
	typeIDMap map[types.BlockType][]string
	idWordMap map[string]*Word
	idLineMap map[string]*Line
}

func newPageParser(bp *blockParser, page *Page) *pageParser {
	typeIDMap := make(map[types.BlockType][]string)
	for k, v := range bp.typeIDMap {
		ids := make([]string, 0)

		for _, id := range v {
			if page.isChild(id) {
				ids = append(ids, id)
			}
		}

		typeIDMap[k] = ids
	}

	return &pageParser{
		bp:        bp,
		page:      page,
		typeIDMap: typeIDMap,
		idWordMap: make(map[string]*Word),
	}
}

func (pp *pageParser) newWord(b types.Block) *Word {
	if val, ok := pp.idWordMap[aws.ToString(b.Id)]; ok {
		return val
	}

	word := &Word{
		base:     newBase(b, pp.page),
		text:     aws.ToString(b.Text),
		textType: b.TextType,
	}

	pp.idWordMap[word.id] = word

	return word
}

func (pp *pageParser) createWords() []*Word {
	words := make([]*Word, 0, len(pp.idWordMap))

	for _, w := range pp.idWordMap {
		if w.line == nil {
			w.line = &Line{
				base: base{
					id:          uuid.New().String(),
					confidence:  w.Confidence(),
					blockType:   types.BlockTypeLine,
					boundingBox: w.BoundingBox(),
					page:        pp.page,
				},
				words: []*Word{w},
			}
		}

		words = append(words, w)
	}

	return words
}

func (pp *pageParser) createLines() []*Line {
	ids := pp.blockTypeIDs(types.BlockTypeLine)
	lines := make([]*Line, 0, len(ids))

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		line := &Line{
			base: newBase(b, pp.page),
		}

		rIDs := filterRelationshipIDsByType(b, types.RelationshipTypeChild)
		words := make([]*Word, 0, len(rIDs))

		for _, rid := range rIDs {
			wb := pp.bp.blockByID(rid)
			word := pp.newWord(wb)
			word.line = line
			words = append(words, word)
		}

		sort.Slice(words, func(i, j int) bool {
			return words[i].BoundingBox().Left() < words[j].BoundingBox().Left() ||
				(words[i].BoundingBox().Left() == words[j].BoundingBox().Left() &&
					words[i].BoundingBox().Top() < words[j].BoundingBox().Top())
		})

		line.words = words
		lines = append(lines, line)
	}

	pp.idLineMap = make(map[string]*Line, len(lines))
	for _, l := range lines {
		pp.idLineMap[l.id] = l
	}

	return lines
}

func (pp *pageParser) createKeyValues() []*KeyValue {
	ids := pp.blockTypeIDs(types.BlockTypeKeyValueSet)
	keyValues := make([]*KeyValue, 0, len(ids))

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		if !slices.Contains(b.EntityTypes, types.EntityTypeKey) {
			continue
		}

		key := &Key{
			base: newBase(b, pp.page),
		}

		for _, wid := range filterRelationshipIDsByType(b, types.RelationshipTypeChild) {
			wb := pp.bp.blockByID(wid)
			word := pp.newWord(wb)
			key.words = append(key.words, word)
		}

		valueIDs := filterRelationshipIDsByType(b, types.RelationshipTypeValue)
		v := pp.bp.blockByID(valueIDs[0])

		value := &Value{
			base: newBase(v, pp.page),
		}

		for _, cid := range filterRelationshipIDsByType(v, types.RelationshipTypeChild) {
			wb := pp.bp.blockByID(cid)
			if wb.BlockType == types.BlockTypeWord {
				word := pp.newWord(wb)
				value.words = append(value.words, word)
			} else if wb.BlockType == types.BlockTypeSelectionElement {
				value.selectionElement = &SelectionElement{
					base:   newBase(wb, pp.page),
					status: wb.SelectionStatus,
				}
			}
		}

		kv := &KeyValue{
			key:   key,
			value: value,
			page:  pp.page,
		}

		keyValues = append(keyValues, kv)

		var (
			added  bool
			delIDs []string
		)

		for _, pl := range pp.page.Layouts() {
			if is := pl.BoundingBox().Intersection(kv.BoundingBox()); is != nil {
				if !added {
					pl.AddChildren(kv)

					added = true
				}

			wordloop:
				for _, w := range kv.Words() {
					pl.children = slices.DeleteFunc(pl.children, func(lc LayoutChild) bool {
						return lc.ID() == w.line.ID()
					})

					if len(pl.children) == 0 {
						break wordloop
					}
				}

				if len(pl.children) == 0 {
					delIDs = append(delIDs, pl.ID())
				} else {
					pl.boundingBox = NewEnclosingBoundingBox(pl.children...)
				}
			}
		}

		pp.page.layouts = slices.DeleteFunc(pp.page.layouts, func(l *Layout) bool {
			return slices.Contains(delIDs, l.ID())
		})
	}

	return keyValues
}

func (pp *pageParser) createLayouts() []*Layout {
	ids := pp.blockTypeIDs(types.BlockType("LAYOUT"))
	layouts := make([]*Layout, 0, len(ids))

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		var layout *Layout
		switch b.BlockType { // nolint exhaustive
		case types.BlockTypeLayoutList:
			layout = &Layout{
				base: newBase(b, pp.page),
			}

			for _, r := range b.Relationships {
				if r.Type == types.RelationshipTypeChild {
					for _, ri := range r.Ids {
						l := pp.bp.blockByID(ri)

						leafLayout := &Layout{
							base:       newBase(l, pp.page),
							noNewLines: true,
						}

						for _, r := range l.Relationships {
							if r.Type == types.RelationshipTypeChild {
								for _, ri := range r.Ids {
									leafLayout.AddChildren(pp.idLineMap[ri])
								}
							}
						}

						layout.AddChildren(leafLayout)
					}
				}
			}
		case types.BlockTypeLayoutText, types.BlockTypeLayoutSectionHeader, types.BlockTypeLayoutTitle:
			layout = &Layout{
				base:       newBase(b, pp.page),
				noNewLines: true,
			}
		default:
			layout = &Layout{
				base:       newBase(b, pp.page),
				noNewLines: false,
			}
		}

		for _, r := range b.Relationships {
			if r.Type == types.RelationshipTypeChild {
				for _, ri := range r.Ids {
					c := pp.bp.blockByID(ri)

					if c.BlockType == types.BlockTypeLine {
						layout.children = append(layout.children, pp.idLineMap[ri])
					} else {
						fmt.Println("TODO LAYOUT", c.BlockType)
					}
				}
			}
		}

		layouts = append(layouts, layout)
	}

	if len(layouts) == 0 {
		layouts = make([]*Layout, 0, len(pp.page.Lines()))

		for _, line := range pp.page.Lines() {
			layout := &Layout{
				base: base{
					id:          uuid.New().String(),
					confidence:  line.Confidence(),
					blockType:   types.BlockTypeLayoutText,
					boundingBox: line.BoundingBox(),
					page:        pp.page,
				},
				noNewLines: false,
			}

			layout.AddChildren(line)

			layouts = append(layouts, layout)
		}
	}

	return layouts
}

func (pp *pageParser) createTables() []*Table {
	ids := pp.blockTypeIDs(types.BlockTypeTable)
	tables := make([]*Table, 0, len(ids))

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		table := &Table{
			base: newBase(b, pp.page),
		}

		for _, cid := range filterRelationshipIDsByType(b, types.RelationshipTypeChild) {
			c := pp.bp.blockByID(cid)

			if c.BlockType == types.BlockTypeCell {
				cell := &TableCell{
					base:        newBase(c, pp.page),
					rowIndex:    int(aws.ToInt32(c.RowIndex)),
					columnIndex: int(aws.ToInt32(c.ColumnIndex)),
					rowSpan:     int(aws.ToInt32(c.RowSpan)),
					columnSpan:  int(aws.ToInt32(c.ColumnSpan)),
					entityTypes: c.EntityTypes,
				}

				for _, rid := range filterRelationshipIDsByType(c, types.RelationshipTypeChild) {
					c := pp.bp.blockByID(rid)

					switch c.BlockType { // nolint exhaustive
					case types.BlockTypeWord:
						word := pp.newWord(c)
						word.tableCell = cell

						cell.words = append(cell.words, word)
					case types.BlockTypeSelectionElement:
						// TODO
						fmt.Println("TODO SelectionElement TABLE CELL")
					}
				}

				table.cells = append(table.cells, cell)
			}
		}

		for _, id := range filterRelationshipIDsByType(b, types.RelationshipTypeTableTitle) {
			t := pp.bp.blockByID(id)

			title := &TableTitle{
				base: newBase(t, pp.page),
			}

			for _, rid := range filterRelationshipIDsByType(t, types.RelationshipTypeChild) {
				w := pp.bp.blockByID(rid)
				if w.BlockType == types.BlockTypeWord {
					word := pp.newWord(w)
					title.words = append(title.words, word)
				}
			}

			table.title = title
		}

		for _, id := range filterRelationshipIDsByType(b, types.RelationshipTypeTableFooter) {
			f := pp.bp.blockByID(id)

			footer := &TableFooter{
				base: newBase(f, pp.page),
			}

			for _, rid := range filterRelationshipIDsByType(f, types.RelationshipTypeChild) {
				w := pp.bp.blockByID(rid)
				if w.BlockType == types.BlockTypeWord {
					footer.words = append(footer.words, pp.newWord(w))
				}
			}

			table.footers = append(table.footers, footer)
		}

		tables = append(tables, table)

		var (
			added  bool
			delIDs []string
		)

		for _, pl := range pp.page.Layouts() {
			if is := pl.BoundingBox().Intersection(table.BoundingBox()); is != nil {
				if !added {
					pl.AddChildren(table)

					added = true
				}

			wordloop:
				for _, w := range table.Words() {
					pl.children = slices.DeleteFunc(pl.children, func(lc LayoutChild) bool {
						return lc.ID() == w.line.ID()
					})

					if len(pl.children) == 0 {
						break wordloop
					}
				}

				if len(pl.children) == 0 {
					delIDs = append(delIDs, pl.ID())
				} else {
					pl.boundingBox = NewEnclosingBoundingBox(pl.children...)
				}
			}
		}

		pp.page.layouts = slices.DeleteFunc(pp.page.layouts, func(l *Layout) bool {
			return slices.Contains(delIDs, l.ID())
		})
	}

	return tables
}

func (pp *pageParser) createQueries() []*Query {
	ids := pp.blockTypeIDs(types.BlockTypeQuery)
	queries := make([]*Query, 0, len(ids))

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		rIDs := filterRelationshipIDsByType(b, types.RelationshipTypeAnswer)

		results := make([]*QueryResult, len(rIDs))

		for i, id := range rIDs {
			rb := pp.bp.blockByID(id)
			results[i] = &QueryResult{
				base: newBase(rb, pp.page),
				text: aws.ToString(rb.Text),
			}
		}

		queries = append(queries, &Query{
			id:         aws.ToString(b.Id),
			text:       aws.ToString(b.Query.Text),
			alias:      aws.ToString(b.Query.Alias),
			queryPages: b.Query.Pages,
			results:    results,
			page:       pp.page,
			raw:        b,
		})
	}

	return queries
}

func (pp *pageParser) createSignatures() []*Signature {
	ids := pp.blockTypeIDs(types.BlockTypeSignature)
	signatures := make([]*Signature, 0, len(ids))

	layouts := pp.page.Layouts()
	sort.Slice(layouts, func(i, j int) bool {
		return layouts[i].BoundingBox().Top() < layouts[j].BoundingBox().Top()
	})

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		signature := &Signature{
			base: newBase(b, pp.page),
		}

		for _, l := range layouts {
			if is := l.BoundingBox().Intersection(signature.BoundingBox()); is != nil {
				if is.Area() > signature.BoundingBox().Area()*threshold {
					l.AddChildren(signature)
					break
				}
			}
		}

		signatures = append(signatures, signature)
	}

	return signatures
}

func (pp *pageParser) blockTypeIDs(blockType types.BlockType) []string {
	return pp.typeIDMap[blockType]
}
