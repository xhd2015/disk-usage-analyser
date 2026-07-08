# Scenario

**Feature**: scan CLI shared fixture harness

```
scan [PATH] [--json] [--threshold] [--max-depth] -> usagescan.Scan / RunCLI -> TreeResult -> text tree or JSON
```

## Preconditions

- Tests use isolated temporary fixture trees under `t.TempDir()`.
- Fixture sizes use exact byte counts so numeric assertions are deterministic.
- The `usagescan` package provides `Scan`, `RunCLI`, `CLIOptions`, `ScanOptions`, `TreeResult`, and `TreeNode`.
- `req.FixtureDir` is always an absolute path suitable as the resolved scan root.
- Default programmatic `ScanOptions` mirror CLI text defaults: threshold `1M`, maxDepth `3`.
- `Threshold: 0` disables display filtering (all nodes shown); used in `basics/` and `sorting/`.

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
- Scan-mode leaves assert numeric `int64` fields primarily; CLI text leaves assert aligned output.

```go
import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"disk-usage-analyser/usagescan"
)

const defaultThreshold = 1 << 20 // 1M

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
		Threshold: defaultThreshold,
		MaxDepth:  3,
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
```