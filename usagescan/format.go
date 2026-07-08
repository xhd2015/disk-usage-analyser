package usagescan

import (
	"fmt"
	"strings"
)

type treeRow struct {
	left string
	size int64
}

func formatTreeText(result TreeResult) string {
	rows := []treeRow{}
	collectTreeRows(result.Tree.Children, "", &rows)

	maxLeft := 0
	for _, row := range rows {
		if n := len(row.left); n > maxLeft {
			maxLeft = n
		}
	}
	sizeCol := maxLeft + 2

	lines := []string{
		fmt.Sprintf("PATH: %s", result.Path),
		fmt.Sprintf("TOTAL: %s", formatHumanSize(result.TotalSize)),
		fmt.Sprintf("THRESHOLD: %s", FormatCompactHumanSize(result.Threshold)),
		fmt.Sprintf("MAX-DEPTH: %d", result.MaxDepth),
		"",
		".",
	}
	for _, row := range rows {
		padding := sizeCol - len(row.left)
		lines = append(lines, row.left+strings.Repeat(" ", padding)+FormatCompactHumanSize(row.size))
	}
	lines = append(lines, "", "")
	return strings.Join(lines, "\n")
}

func collectTreeRows(children []TreeNode, prefix string, rows *[]treeRow) {
	for i, child := range children {
		isLast := i == len(children)-1
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		name := child.Name
		if child.IsDir {
			name += "/"
		}

		left := prefix + connector + name
		*rows = append(*rows, treeRow{left: left, size: child.Size})

		if len(child.Children) > 0 {
			childPrefix := prefix + "│   "
			if isLast {
				childPrefix = prefix + "    "
			}
			collectTreeRows(child.Children, childPrefix, rows)
		}
	}
}

func formatHumanSize(size int64) string {
	if size <= 1024 {
		return fmt.Sprintf("%d B", size)
	}
	units := []string{"K", "M", "G", "T", "P"}
	value := float64(size)
	for _, unit := range units {
		value /= 1024
		if value < 1024 {
			if value == float64(int64(value)) {
				return fmt.Sprintf("%.0f%s", value, unit)
			}
			return fmt.Sprintf("%.1f%s", value, unit)
		}
	}
	return fmt.Sprintf("%.1fE", value/1024)
}