package usagescan

import (
	"fmt"
)

type treeRow struct {
	left string
	size int64
}

func formatTreeText(result TreeResult) string {
	// Tree-only live/capture path (no SOURCE, no TOP).
	view := ViewResult{
		ScanPath:  result.Path,
		TotalSize: result.TotalSize,
		Min:       result.Min,
		MaxDepth:  result.MaxDepth,
		Tree:      result.Tree,
	}
	return FormatViewText(view, ViewOptions{})
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
