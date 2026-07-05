package nmpipeline

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"disk-usage-analyser/nminventory"
)

func formatSignedSize(size int64) string {
	if size == 0 {
		return "0B"
	}
	sign := "+"
	if size < 0 {
		sign = "-"
	}
	return sign + nminventory.FormatCompactHumanSize(abs64(size))
}

func abs64(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func truncatePath(path string, max int) string {
	if len(path) <= max {
		return path
	}
	if max <= 1 {
		return path[:max]
	}
	return "…" + path[len(path)-(max-1):]
}

// FormatTable writes an aligned CLI table for pipeline rows.
func FormatTable(w io.Writer, rows []Row) error {
	headers := []string{
		"path",
		"before_size",
		"before_pnpm",
		"before_bun",
		"after_size",
		"after_pnpm",
		"after_bun",
		"shared_size_added",
	}

	cells := make([][]string, 0, len(rows)+1)
	cells = append(cells, headers)

	for _, row := range rows {
		cells = append(cells, []string{
			truncatePath(displayPath(row.Entry.Path), 72),
			nminventory.FormatCompactHumanSize(row.BeforeTotal),
			nminventory.FormatCompactHumanSize(row.BeforePnpm),
			nminventory.FormatCompactHumanSize(row.BeforeBun),
			nminventory.FormatCompactHumanSize(row.AfterTotal),
			nminventory.FormatCompactHumanSize(row.AfterPnpm),
			nminventory.FormatCompactHumanSize(row.AfterBun),
			formatSignedSize(row.SharedAdded()),
		})
	}

	widths := columnWidths(cells)
	var b strings.Builder
	writeRow(&b, cells[0], widths, false)
	for _, row := range cells[1:] {
		writeRow(&b, row, widths, true)
	}

	var totalAdded int64
	for _, row := range rows {
		totalAdded += row.SharedAdded()
	}
	fmt.Fprintf(&b, "\nTOTAL shared_size_added: %s across %d path(s)\n", formatSignedSize(totalAdded), len(rows))

	_, err := io.WriteString(w, b.String())
	return err
}

func columnWidths(rows [][]string) []int {
	if len(rows) == 0 {
		return nil
	}
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	return widths
}

func writeRow(b *strings.Builder, cells []string, widths []int, numeric bool) {
	for i, cell := range cells {
		if i == 0 {
			b.WriteString(padRight(cell, widths[i]))
		} else {
			b.WriteString(padLeft(cell, widths[i]))
		}
		if i < len(cells)-1 {
			b.WriteByte(' ')
		}
	}
	b.WriteByte('\n')
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func sortRowsBySharedAdded(rows []Row) {
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].SharedAdded() > rows[j].SharedAdded()
	})
}