# Truncate Path Keep Suffix — Pure Logic Tests

Unit tests for `truncatePathKeepSuffix` in `pathDisplay.ts`: prefix-hide truncation
that keeps the path suffix fully visible. Logic runs in Node via `tsx` (no browser).

## Version

0.0.2

# DSN (Domain Specific Notion)

**TruncatedPath** in TmpFilesAnalyse receives a filesystem path and column width budget.
When the path exceeds the visible character limit, **truncatePathKeepSuffix** hides the
prefix with an ellipsis character and returns a display string whose suffix (e.g.
`/my-repo/node_modules`) remains intact. When the path fits, the helper returns the
path unchanged. The helper prefers breaking at `/` boundaries when choosing where to
slice the prefix.

## Test Tree

```
truncate-path-keep-suffix/
├── fits-within-limit/
├── truncates-with-prefix-ellipsis/
└── prefers-slash-boundary/
```

## Test Index

| Leaf | Op | Fixture |
|------|-----|---------|
| fits-within-limit | truncate-path | testdata/fixture.json |
| truncates-with-prefix-ellipsis | truncate-path | testdata/fixture.json |
| prefers-slash-boundary | truncate-path | testdata/fixture.json |

## How to Run

```sh
doctest vet ./tests/tmp-analyse-frontend-test-cases/named-section/path-truncation/truncate-path-keep-suffix
doctest test ./tests/tmp-analyse-frontend-test-cases/named-section/path-truncation/truncate-path-keep-suffix
```

Requires **Node.js** and network for first `npx tsx` fetch. The harness imports
`disk-usage-analyser-react/src/pathDisplay.ts` (RED until implemented).

```go
import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type Request struct {
	Op          string
	FixtureFile string
}

type Response struct {
	Output string
	JSON   map[string]interface{}
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Op == "" {
		t.Fatal("req.Op is required")
	}
	if req.FixtureFile == "" {
		t.Fatal("req.FixtureFile is required")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fixturePath := filepath.Join(cwd, req.FixtureFile)
	if _, err := os.Stat(fixturePath); err != nil {
		return nil, fmt.Errorf("fixture %s: %w", fixturePath, err)
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	harness := filepath.Join(projectRoot, "tests/tmp-analyse-frontend-test-cases/named-section/path-truncation/truncate-path-keep-suffix/path-display-harness.ts")
	cmd := exec.Command("npx", "--yes", "tsx", harness, "--op", req.Op, "--fixture", fixturePath)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return &Response{Output: output}, fmt.Errorf("path display harness: %w\n%s", err, output)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(out, &parsed); err != nil {
		return &Response{Output: output}, fmt.Errorf("parse harness JSON: %w\n%s", err, output)
	}

	return &Response{Output: output, JSON: parsed}, nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found")
		}
		dir = parent
	}
}
```