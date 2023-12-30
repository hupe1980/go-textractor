package textractor

import (
	"fmt"
	"slices"
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/textract/types"
	"github.com/hupe1980/go-textractor/internal"
)

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
					id:          "", // TODO
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
		words := make([]*Word, len(rIDs))

		for i, rid := range rIDs {
			wb := pp.bp.blockByID(rid)
			word := pp.newWord(wb)
			word.line = line
			words[i] = word
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
	}

	return keyValues
}

func (pp *pageParser) createLayouts() []Layout {
	ids := pp.blockTypeIDs(types.BlockType("LAYOUT"))
	layouts := make([]Layout, 0, len(ids))

	readingOrder := 0

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		if b.BlockType == types.BlockTypeLayoutList {
			containerLayout := &ContainerLayout{
				layout: newLayout(b, pp.page, readingOrder),
			}

			for _, r := range b.Relationships {
				if r.Type == types.RelationshipTypeChild {
					for i, ri := range r.Ids {
						leaf := pp.bp.blockByID(ri)
						leafLayout := &LeafLayout{
							layout:     newLayout(leaf, pp.page, i),
							noNewLines: true,
						}

						for _, r := range b.Relationships {
							if r.Type == types.RelationshipTypeChild {
								for _, ri := range r.Ids {
									leafLayout.lines = append(leafLayout.lines, pp.idLineMap[ri])
								}
							}
						}

						containerLayout.layouts = append(containerLayout.layouts, leafLayout)
					}
				}
			}

			layouts = append(layouts, containerLayout)
		} else {
			var layout *LeafLayout
			switch b.BlockType { // nolint exhaustive
			case types.BlockTypeLayoutText, types.BlockTypeLayoutSectionHeader, types.BlockTypeLayoutTitle:
				layout = &LeafLayout{
					layout:     newLayout(b, pp.page, readingOrder),
					noNewLines: true,
				}
			default:
				layout = &LeafLayout{
					layout:     newLayout(b, pp.page, readingOrder),
					noNewLines: false,
				}
			}

			for _, r := range b.Relationships {
				if r.Type == types.RelationshipTypeChild {
					for _, ri := range r.Ids {
						c := pp.bp.blockByID(ri)

						if c.BlockType == types.BlockTypeLine {
							layout.lines = append(layout.lines, pp.idLineMap[ri])
						} else {
							fmt.Println("TODO LAYOUT", readingOrder, c.BlockType)
						}
					}
				}
			}

			layouts = append(layouts, layout)
		}

		readingOrder++
	}

	if len(layouts) == 0 {
		lines := internal.Values(pp.idLineMap)
		layouts = make([]Layout, 0, len(lines))

		for i, line := range lines {
			layout := &LeafLayout{
				layout:     newLayout(line.Raw(), pp.page, i),
				noNewLines: false,
			}

			layout.lines = append(layout.lines, line)
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
			fmt.Println("TODO TABLE TITLE", id)
		}

		for _, id := range filterRelationshipIDsByType(b, types.RelationshipTypeTableFooter) {
			fmt.Println("TODO TABLE FOOTER", id)
		}

		tables = append(tables, table)
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

	for _, id := range ids {
		b := pp.bp.blockByID(id)

		signatures = append(signatures, &Signature{
			base: newBase(b, pp.page),
		})
	}

	return signatures
}

func (pp *pageParser) blockTypeIDs(blockType types.BlockType) []string {
	return pp.typeIDMap[blockType]
}