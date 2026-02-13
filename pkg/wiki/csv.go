package wiki

import (
	"io"
	"strings"
)

func cellCSV(cell string) string {
	csvText := string(cell)

	if strings.Contains(csvText, `"`) {
		csvText = strings.ReplaceAll(csvText, `"`, `""`)
	}

	if strings.ContainsAny(csvText, "\r\n\",") {
		csvText = `"` + csvText + `"`
	}

	return csvText
}

func WriteCSV(table Table, f io.WriteCloser) error {
	defer f.Close()

	for _, row := range table.Rows {
		cells := make([]string, len(row.Cells))

		for i, c := range row.Cells {
			cells[i] = cellCSV(c)
		}

		_, err := io.WriteString(f, strings.Join(cells, ",")+"\n")

		if err != nil {
			return err
		}
	}
	return nil
}
