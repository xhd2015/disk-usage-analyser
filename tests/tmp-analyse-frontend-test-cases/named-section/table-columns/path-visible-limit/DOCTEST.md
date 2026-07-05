# Path Visible Limit — Unit Test

Verifies `PATH_VISIBLE_CHAR_LIMIT` in `pathDisplay.ts` is **56** so node_modules
path cells show more suffix before ellipsis truncation.

## Version

0.0.2

# DSN (Domain Specific Notion)

**TruncatedPath** in TmpFilesAnalyse reads `PATH_VISIBLE_CHAR_LIMIT` from
**pathDisplay.ts** when calling **truncatePathKeepSuffix**. A higher limit keeps
more of the filesystem path visible in the grid's `1fr` Path column before the
prefix is hidden with an ellipsis.

## Test Tree

```
path-visible-limit/
└── (leaf)
```

## Test Index

| Leaf | Op |
|------|-----|
| path-visible-limit | path-visible-limit |

## How to Run

```sh
doctest vet ./tests/tmp-analyse-frontend-test-cases/named-section/table-columns/path-visible-limit
doctest test ./tests/tmp-analyse-frontend-test-cases/named-section/table-columns/path-visible-limit
```

Requires **Node.js** and network for first `npx tsx` fetch.

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
	Op string
}

type Response struct {
	Output string
	JSON   map[string]interface{}
}

func Run(t *testing.T, req *Request) (*Response, error) {
	if req.Op == "" {
		req.Op = "path-visible-limit"
	}

	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	harness := filepath.Join(projectRoot, "tests/tmp-analyse-frontend-test-cases/named-section/table-columns/path-visible-limit/path-visible-limit-harness.ts")
	cmd := exec.Command("npx", "--yes", "tsx", harness, "--op", req.Op)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return &Response{Output: output}, fmt.Errorf("path visible limit harness: %w\n%s", err, output)
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