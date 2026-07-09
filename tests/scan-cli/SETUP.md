# Scenario

**Feature**: scan CLI two-phase harness (live FS and inspect JSON)

```
# Phase 1: TreeSource
scan [PATH] [--min] [--max-depth] -> LiveTreeSource -> TreeResult
scan --inspect FILE [...]         -> JSONTreeSource -> TreeResult

# Phase 2: View
TreeResult + ViewOptions -> tree section [+ match section] | JSON capture/ViewResult
```

## Preconditions

- Tests use isolated temporary fixture trees under `t.TempDir()`.
- Fixture sizes use exact byte counts so numeric assertions are deterministic.
- The `usagescan` package provides `Scan`, `RunCLI`, `CLIOptions`, `ScanOptions`, `TreeResult`,
  and `TreeNode` with field **`Min`** (JSON `min`; replaces `Threshold` / `threshold`).
- `req.FixtureDir` is always an absolute path suitable as the resolved live scan root.
- Default programmatic `ScanOptions` mirror CLI live text defaults: **min** `1M`, **maxDepth** `3`.
- `Min: 0` disables display filtering (all nodes shown); used in `basics/` and `sorting/`.
- Inspect fixtures write fixed TreeResult JSON with field **`min`** only (no `threshold`).

## Steps

1. Create a fresh fixture root for each leaf.
2. Set `req.FixtureDir` to the fixture path (or `req.Args` for CLI leaves).
3. Run `usagescan.Scan`, `usagescan.RunCLI`, or `run.RunWithOptions` via the root `Run` harness.

## Context

- Default `req.Mode` is `scan` (structured `TreeResult` assertions).
- CLI leaves set `req.Mode` to `cli` or `dispatch`.
- `RunCLI` receives args **after** the `scan` token (no `scan` prefix in `req.Args`).
- Text tree lines place the name after `├──`/`└──` and right-pad to a shared size column
  (`maxLen(prefix+connector+name)` + at least two spaces before `FormatCompactHumanSize`).
- Summary label is **`MIN:`** (not `THRESHOLD:`). Inspect views add **`SOURCE:`**.
- Scan-mode leaves assert numeric `int64` fields primarily; CLI text leaves assert aligned output.
- User-facing stdout ends with a trailing blank line after the last content line.

```go
import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"disk-usage-analyser/usagescan"
)

const defaultMin = 1 << 20 // 1M

func Setup(t *testing.T, req *Request) error {
	req.FixtureDir = filepath.Join(t.TempDir(), "fixture")
	if err := os.MkdirAll(req.FixtureDir, 0755); err != nil {
		return err
	}
	if req.Mode == "" {
		req.Mode = "scan"
	}
	return nil
}

func defaultScanOpts() usagescan.ScanOptions {
	return usagescan.ScanOptions{
		Min:      defaultMin,
		MaxDepth: 3,
	}
}

func mkdir(t *testing.T, base string, rel string) string {
	t.Helper()
	dir := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("mkdir %s: %v", rel, err)
	}
	return dir
}

func writeSizedFile(t *testing.T, base string, rel string, size int64) string {
	t.Helper()
	path := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatalf("mkdir parent %s: %v", rel, err)
	}
	data := make([]byte, size)
	for i := range data {
		data[i] = 'x'
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatalf("write %s: %v", rel, err)
	}
	return path
}

func treeChildByName(children []usagescan.TreeNode, name string) *usagescan.TreeNode {
	for i := range children {
		if children[i].Name == name {
			return &children[i]
		}
	}
	return nil
}

func assertTreeChildrenSorted(t *testing.T, children []usagescan.TreeNode) {
	t.Helper()
	for i := 1; i < len(children); i++ {
		prev, cur := children[i-1], children[i]
		if cur.Size > prev.Size {
			t.Fatalf("children not sorted by size desc at %d: %#v before %#v", i, prev, cur)
		}
		if cur.Size == prev.Size && !prev.IsDir && cur.IsDir {
			t.Fatalf("dirs must sort before files when sizes tie: %#v before %#v", prev, cur)
		}
	}
}

func assertRootTree(t *testing.T, result *usagescan.TreeResult, fixtureDir string) *usagescan.TreeNode {
	t.Helper()
	if result.Path != fixtureDir {
		t.Fatalf("path: expected %q, got %q", fixtureDir, result.Path)
	}
	if result.Tree.Name != "." {
		t.Fatalf("root tree name: expected %q, got %q", ".", result.Tree.Name)
	}
	if result.Tree.Path != fixtureDir {
		t.Fatalf("root tree path: expected %q, got %q", fixtureDir, result.Tree.Path)
	}
	if !result.Tree.IsDir {
		t.Fatal("root tree must be a directory")
	}
	if result.Tree.Depth != 0 {
		t.Fatalf("root depth: expected 0, got %d", result.Tree.Depth)
	}
	return &result.Tree
}

func stdoutEndsWithBlankLine(t *testing.T, stdout string) {
	t.Helper()
	if stdout == "" {
		t.Fatal("stdout is empty")
	}
	if !strings.HasSuffix(stdout, "\n\n") {
		t.Fatalf("stdout must end with trailing blank line after last content line; got %q", stdout)
	}
}

func assertTreeSizeColumnAligned(t *testing.T, stdout string) {
	t.Helper()
	sizePattern := regexp.MustCompile(`\d+(?:\.\d+)?[KMGTPE]?B$`)
	var columns []int
	for _, line := range strings.Split(stdout, "\n") {
		if !strings.Contains(line, "──") {
			continue
		}
		loc := sizePattern.FindStringIndex(line)
		if loc == nil {
			continue
		}
		if loc[0] < 2 || line[loc[0]-1] != ' ' || line[loc[0]-2] != ' ' {
			t.Fatalf("tree line must have at least two spaces before size: %q", line)
		}
		columns = append(columns, loc[0])
	}
	if len(columns) < 2 {
		t.Fatalf("expected at least two tree lines with sizes, got %d in:\n%s", len(columns), stdout)
	}
	first := columns[0]
	for i, col := range columns[1:] {
		if col != first {
			t.Fatalf("size column misaligned: line %d starts at %d, expected %d:\n%s", i+2, col, first, stdout)
		}
	}
}

// writeTreeResultJSON writes a TreeResult-shaped JSON file with field "min" (post-rename).
// scanRoot is used for path fields; the file is written to jsonPath.
func writeTreeResultJSON(t *testing.T, jsonPath, scanRoot string, totalSize int64, min int64, maxDepth int, tree usagescan.TreeNode) {
	t.Helper()
	payload := map[string]any{
		"path":      scanRoot,
		"totalSize": totalSize,
		"min":       min,
		"maxDepth":  maxDepth,
		"tree":      tree,
	}
	raw, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatalf("marshal TreeResult JSON: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0755); err != nil {
		t.Fatalf("mkdir for JSON: %v", err)
	}
	if err := os.WriteFile(jsonPath, raw, 0644); err != nil {
		t.Fatalf("write JSON %s: %v", jsonPath, err)
	}
}

// sampleInspectTree builds a depth≥2 tree under scanRoot for inspect leaves.
// Layout (size-desc children): huge.bin 400, mid.bin 200, big/(deep.bin 50);
// total 650. Top-2 non-root ranking is uniquely huge then mid (big/deep are smaller).
func sampleInspectTree(scanRoot string) (usagescan.TreeNode, int64) {
	deep := usagescan.TreeNode{
		Name:  "deep.bin",
		Path:  filepath.Join(scanRoot, "big", "deep.bin"),
		Size:  50,
		IsDir: false,
		Depth: 2,
	}
	big := usagescan.TreeNode{
		Name:     "big",
		Path:     filepath.Join(scanRoot, "big"),
		Size:     50,
		IsDir:    true,
		Depth:    1,
		Children: []usagescan.TreeNode{deep},
	}
	mid := usagescan.TreeNode{
		Name:  "mid.bin",
		Path:  filepath.Join(scanRoot, "mid.bin"),
		Size:  200,
		IsDir: false,
		Depth: 1,
	}
	huge := usagescan.TreeNode{
		Name:  "huge.bin",
		Path:  filepath.Join(scanRoot, "huge.bin"),
		Size:  400,
		IsDir: false,
		Depth: 1,
	}
	root := usagescan.TreeNode{
		Name:     ".",
		Path:     scanRoot,
		Size:     650,
		IsDir:    true,
		Depth:    0,
		// size desc: huge, mid, big
		Children: []usagescan.TreeNode{huge, mid, big},
	}
	return root, 650
}

func firstJSONObjectLine(t *testing.T, stdout string) string {
	t.Helper()
	content := strings.TrimRight(stdout, "\n")
	lines := strings.Split(content, "\n")
	if len(lines) != 1 {
		t.Fatalf("expected one JSON line, got %d lines:\n%s", len(lines), stdout)
	}
	return lines[0]
}

func stdoutHasNoTopSection(t *testing.T, stdout string) {
	t.Helper()
	for _, line := range strings.Split(stdout, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "TOP ") || trim == "TOP" {
			t.Fatalf("expected no TOP match section, got line %q in:\n%s", line, stdout)
		}
	}
}
```
