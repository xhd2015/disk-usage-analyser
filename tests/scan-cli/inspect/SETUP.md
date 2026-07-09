# Scenario

**Feature**: Phase 1 JSONTreeSource via `scan --inspect FILE`

```
# load saved TreeResult (field min), then shared View
scan --inspect FILE [view flags] -> JSONTreeSource.Load -> View(tree [+ matches])
```

## Preconditions

- Leaves write fixed TreeResult JSON with **`min`** (never `threshold`).
- Default inspect view: **maxDepth 1**, **min 0**.
- `SOURCE:` appears in human text; matches use full loaded tree for ranking.

## Context

- Option B: `--top` / `--find` / `--suffix` add a match section after the tree.
- `--at` alone focuses the tree only (no TOP section).
- `--json` without pure live capture emits ViewResult (`scanPath`, `sourceFile`, `tree`, optional `matches`).

```go
import (
	"os"
	"path/filepath"
)

// inspectPaths is filled by prepareSampleInspect for leaf asserts.
type inspectPaths struct {
	ScanRoot string
	JSONPath string
}

func Setup(t *testing.T, req *Request) error {
	req.Mode = "cli"
	return nil
}

func prepareSampleInspect(t *testing.T, req *Request) inspectPaths {
	t.Helper()
	scanRoot := filepath.Join(t.TempDir(), "scan-root")
	if err := os.MkdirAll(scanRoot, 0755); err != nil {
		t.Fatalf("mkdir scan-root: %v", err)
	}
	jsonPath := filepath.Join(t.TempDir(), "tree.json")
	tree, total := sampleInspectTree(scanRoot)
	writeTreeResultJSON(t, jsonPath, scanRoot, total, 0, 24, tree)
	req.FixtureDir = scanRoot
	return inspectPaths{ScanRoot: scanRoot, JSONPath: jsonPath}
}
```
