package wiki

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type Row struct {
	Cells []string
}

type Table struct {
	Rows []Row
}

func hasExcludedClassName(excluded []string, classAttr string) bool {
	for _, ex := range excluded {
		if strings.Contains(classAttr, ex) {
			return true
		}
	}
	return false
}

var forbiddenNodeNames = map[string]bool{
	"style": true,
	"link":  true,
}

func parseNodeRecursive(n *html.Node, includeLineBreaks bool, excludedCSSClassNames []string) []string {
	var result []string

	if n.Type == html.TextNode {
		result = append(result, n.Data)
		return result
	}

	if n.Type == html.ElementNode {
		nodeName := strings.ToLower(n.Data) // TODO: weird

		if (nodeName == "br" && !includeLineBreaks) || forbiddenNodeNames[nodeName] {
			return result
		}

		var classAttr string
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				classAttr = attr.Val
			}
		}

		if hasExcludedClassName(excludedCSSClassNames, classAttr) {
			return result
		}

		if nodeName == "br" {
			result = append(result, "\n")
			return result
		} else {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				result = append(result, parseNodeRecursive(c, includeLineBreaks, excludedCSSClassNames)...)
			}
		}
	}

	return result
}

func parseCell(cell *goquery.Selection, trimCells bool, includeLineBreaks bool, excludedCSSClassNames []string) string {
	node := cell.Get(0)

	var textParts []string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		textParts = append(textParts, parseNodeRecursive(c, includeLineBreaks, excludedCSSClassNames)...)
	}

	text := strings.Join(textParts, "")
	if trimCells {
		text = strings.TrimSpace(text)
	}

	return text
}

func Parse(wikiData WikiResponse, tableSelector string, cssClassNamesToExclude []string, noTrimCells bool, noIncludeLineBreaks bool) ([]Table, error) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(wikiData.Parse.Text.Content))

	if err != nil {
		return nil, err
	}

	var results []*goquery.Selection
	doc.Find(tableSelector).Each(func(i int, s *goquery.Selection) {
		if err == nil {
			results = append(results, s)
		}
	})

	tables := make([]Table, len(results))

	for i, t := range results {
		res := Table{
			Rows: make([]Row, 0),
		}

		t.Find("tr").Each(func(y int, rowSel *goquery.Selection) {
			rowSel.Find("th, td").Each(func(x int, cellSel *goquery.Selection) {

				rowspan := 1
				colspan := 1

				if v, exists := cellSel.Attr("rowspan"); exists {
					fmt.Sscanf(v, "%d", &rowspan)
				}
				if v, exists := cellSel.Attr("colspan"); exists {
					fmt.Sscanf(v, "%d", &colspan)
				}

				for len(res.Rows) <= y {
					res.Rows = append(res.Rows, Row{make([]string, 0)})
				}

				for len(res.Rows[y].Cells) > x && res.Rows[y].Cells[x] != "" {
					x++
				}

				cell := parseCell(cellSel, !noTrimCells, !noIncludeLineBreaks, cssClassNamesToExclude)

				for yy := y; yy < y+rowspan; yy++ {
					for len(res.Rows) <= yy {
						res.Rows = append(res.Rows, Row{make([]string, 0)})
					}

					row := res.Rows[yy]

					for len(row.Cells) < x+colspan {
						row.Cells = append(row.Cells, "")
					}

					for j := 0; j < colspan; j++ {
						row.Cells[x+j] = cell
					}

					res.Rows[yy] = row
				}
			})
		})

		tables[i] = res
	}

	return tables, nil
}
