# Filter Named Repos By Column — Pure Logic Tests

Unit tests for `filterNamedReposByColumnFilters` in `repositoryScansLayout.ts`: client-side
Git, package.json, and package-manager filters on node_modules rows. Logic runs in Node via
`tsx` (no browser).

## Version

0.0.2

# DSN (Domain Specific Notion)

**TmpFilesAnalyse** accumulates `NamedHit` rows per repo from SSE `named` events. After the
existing **filterNamedRepos** size gate (`showNamedUnder1M`), **filterNamedReposByColumnFilters**
narrows visible hits using tri-state Git and package.json filters plus a package-manager
selector. PM matching uses `hit.packageManager || 'unknown'`. Repos whose hits are all
filtered out disappear from the grouped tree, matching the size-filter empty-repo rule.

## Test Tree

```
filter-named-repos-by-column/
├── defaults-show-all/
├── git-yes/
├── git-no/
├── package-json-yes/
├── package-json-no/
├── pm-pnpm/
├── pm-unknown/
└── combined-filters/
```

## Test Index

| Leaf | Op | Fixture |
|------|-----|---------|
| defaults-show-all | filter-named-column-filters | testdata/fixture.json |
| git-yes | filter-named-column-filters | testdata/fixture.json |
| git-no | filter-named-column-filters | testdata/fixture.json |
| package-json-yes | filter-named-column-filters | testdata/fixture.json |
| package-json-no | filter-named-column-filters | testdata/fixture.json |
| pm-pnpm | filter-named-column-filters | testdata/fixture.json |
| pm-unknown | filter-named-column-filters | testdata/fixture.json |
| combined-filters | filter-named-column-filters | testdata/fixture.json |

## How to Run

```sh
doctest vet ./tests/tmp-analyse-frontend-test-cases/named-section/column-filters/filter-named-repos-by-column
doctest test ./tests/tmp-analyse-frontend-test-cases/named-section/column-filters/filter-named-repos-by-column
```

Requires **Node.js** and network for first `npx tsx` fetch. The harness imports
`disk-usage-analyser-react/src/repositoryScansLayout.ts` (RED until `filterNamedReposByColumnFilters` exists).

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

	harness := filepath.Join(projectRoot, "tests/tmp-analyse-frontend-test-cases/named-section/column-filters/filter-named-repos-by-column/column-filters-harness.ts")
	cmd := exec.Command("npx", "--yes", "tsx", harness, "--op", req.Op, "--fixture", fixturePath)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return &Response{Output: output}, fmt.Errorf("column filters harness: %w\n%s", err, output)
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