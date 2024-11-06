package tableui

import (
	"fmt"
	"io"
	"strings"
)

type Row []string

type TableUI struct {
	headers []string
	rows    []Row
}

func New(headers []string, rows []Row) *TableUI {
	return &TableUI{
		headers: headers,
		rows:    rows,
	}
}

func (t *TableUI) Render(w io.Writer) error {
	for _, row := range t.rows {
		if len(row) != len(t.headers) {
			return fmt.Errorf("row length does not match header length")
		}
	}

	maxSizes := make([]int, len(t.headers))
	for _, row := range t.rows {
		for i, column := range row {
			if len(column) > maxSizes[i] {
				maxSizes[i] = len(column)
			}
		}
	}
	for i, header := range t.headers {
		if len(header) > maxSizes[i] {
			maxSizes[i] = len(header)
		}
	}

	for i, header := range t.headers {
		fmt.Fprintf(w, "%s%s|", strings.Repeat(" ", maxSizes[i]-len(header)), header)
	}

	fmt.Fprintf(w, "\n")
	for i := range t.headers {
		fmt.Fprintf(w, "%s", strings.Repeat("-", maxSizes[i]+1))
	}

	for _, row := range t.rows {
		fmt.Fprintf(w, "\n")
		for i, column := range row {
			fmt.Fprintf(w, "%s%s|", strings.Repeat(" ", maxSizes[i]-len(column)), column)
		}
	}

	fmt.Fprintf(w, "\n")

	return nil
}
