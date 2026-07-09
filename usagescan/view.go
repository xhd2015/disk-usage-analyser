package usagescan

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

// WantMatchSection reports whether Option B match ranking is active.
func WantMatchSection(opts ViewOptions) bool {
	return opts.TopSet || opts.Find != "" || opts.Suffix != ""
}

// BuildView applies phase-2 view options to a loaded TreeResult.
func BuildView(result TreeResult, opts ViewOptions) ViewResult {
	treeRoot := result.Tree
	if opts.AtPath != "" {
		found := false
		if n := FindNodeByPath(result.Tree, opts.AtPath); n != nil {
			treeRoot = *n
			found = true
		} else {
			// fallback: suffix path match
			want := filepath.Clean(opts.AtPath)
			for _, n := range FlattenTree(result.Tree) {
				if strings.HasSuffix(filepath.Clean(n.Path), want) {
					treeRoot = n
					found = true
					break
				}
			}
		}
		if !found {
			treeRoot = TreeNode{
				Name:  ".",
				Path:  result.Path,
				IsDir: true,
				Depth: 0,
			}
		}
	}

	viewTree := pruneTreeRelative(treeRoot, opts.Min, opts.MaxDepth, 0)
	// Always present root as "." for text tree formatting consistency.
	viewTree.Name = "."

	out := ViewResult{
		ScanPath:   result.Path,
		TotalSize:  result.TotalSize,
		Min:        opts.Min,
		MaxDepth:   opts.MaxDepth,
		SourceFile: opts.SourceFile,
		Tree:       viewTree,
	}

	if WantMatchSection(opts) {
		out.Matches = rankMatches(result.Tree, opts)
	}
	return out
}

// pruneTreeRelative filters by min and limits expansion by depth from the view root.
func pruneTreeRelative(n TreeNode, min int64, maxDepth int, depthFromRoot int) TreeNode {
	out := n
	out.Depth = depthFromRoot
	out.Children = nil

	canExpand := maxDepth == 0 || depthFromRoot < maxDepth
	if !canExpand {
		return out
	}

	for _, c := range n.Children {
		if !passesMin(c.Size, min) {
			continue
		}
		child := pruneTreeRelative(c, min, maxDepth, depthFromRoot+1)
		out.Children = append(out.Children, child)
	}
	return out
}

func rankMatches(root TreeNode, opts ViewOptions) []Match {
	nodes := FlattenTree(root)
	topN := opts.Top
	if topN <= 0 {
		topN = 20
	}

	var matches []Match
	for _, n := range nodes {
		if !opts.IncludeRoot && n.Depth == 0 {
			continue
		}
		if !passesMin(n.Size, opts.Min) {
			continue
		}
		if !matchesFind(n, opts.Find) || !matchesSuffix(n, opts.Suffix) {
			continue
		}
		matches = append(matches, matchFromNode(n))
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Size != matches[j].Size {
			return matches[i].Size > matches[j].Size
		}
		return matches[i].Path < matches[j].Path
	})
	if len(matches) > topN {
		matches = matches[:topN]
	}
	return matches
}

func matchFromNode(n TreeNode) Match {
	return Match{
		Name:  n.Name,
		Path:  n.Path,
		Size:  n.Size,
		IsDir: n.IsDir,
		Depth: n.Depth,
	}
}

// FormatViewText renders Option B human output: summary, tree, optional TOP section.
func FormatViewText(view ViewResult, opts ViewOptions) string {
	var b strings.Builder
	fmt.Fprintf(&b, "PATH: %s\n", view.ScanPath)
	fmt.Fprintf(&b, "TOTAL: %s\n", formatHumanSize(view.TotalSize))
	fmt.Fprintf(&b, "MIN: %s\n", FormatCompactHumanSize(view.Min))
	fmt.Fprintf(&b, "MAX-DEPTH: %d\n", view.MaxDepth)
	if view.SourceFile != "" {
		fmt.Fprintf(&b, "SOURCE: %s\n", view.SourceFile)
	}
	b.WriteByte('\n')
	b.WriteString(formatTreeBody(view.Tree))

	if WantMatchSection(opts) {
		topN := opts.Top
		if topN <= 0 {
			topN = 20
		}
		// Header uses the requested/default cap, not the actual match count.
		fmt.Fprintf(&b, "\nTOP %d\n", topN)
		b.WriteString(formatMatchLines(view.Matches))
	}

	// Trailing blank line after last content line.
	b.WriteByte('\n')
	return b.String()
}

func formatTreeBody(root TreeNode) string {
	rows := []treeRow{}
	collectTreeRows(root.Children, "", &rows)

	maxLeft := 0
	for _, row := range rows {
		if n := len(row.left); n > maxLeft {
			maxLeft = n
		}
	}
	sizeCol := maxLeft + 2

	var b strings.Builder
	b.WriteString(".\n")
	for _, row := range rows {
		padding := sizeCol - len(row.left)
		if padding < 2 {
			padding = 2
		}
		b.WriteString(row.left)
		b.WriteString(strings.Repeat(" ", padding))
		b.WriteString(FormatCompactHumanSize(row.size))
		b.WriteByte('\n')
	}
	return b.String()
}

func formatMatchLines(matches []Match) string {
	if len(matches) == 0 {
		return ""
	}
	type row struct {
		size  string
		kind  string
		depth string
		path  string
	}
	rows := make([]row, 0, len(matches))
	maxSize, maxKind, maxDepth := 0, 0, 0
	for _, m := range matches {
		kind := "file"
		if m.IsDir {
			kind = "dir"
		}
		r := row{
			size:  FormatCompactHumanSize(m.Size),
			kind:  kind,
			depth: fmt.Sprintf("d=%d", m.Depth),
			path:  m.Path,
		}
		rows = append(rows, r)
		if len(r.size) > maxSize {
			maxSize = len(r.size)
		}
		if len(r.kind) > maxKind {
			maxKind = len(r.kind)
		}
		if len(r.depth) > maxDepth {
			maxDepth = len(r.depth)
		}
	}
	var b strings.Builder
	for _, r := range rows {
		fmt.Fprintf(&b, "%-*s  %-*s  %-*s  %s\n",
			maxSize, r.size,
			maxKind, r.kind,
			maxDepth, r.depth,
			r.path,
		)
	}
	return b.String()
}
