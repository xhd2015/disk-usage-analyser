# Repository Scans Layout — Pure Logic Tests

Unit tests for `repositoryScansLayout.ts`: size DESC sorting and MiB threshold
filtering for Git Worktrees and Binary files sections. Logic runs in Node via
`tsx` against the extracted TypeScript module (no browser).

## Version

0.0.2

# DSN (Domain Specific Notion)

**TmpFilesAnalyse** receives SSE `repo`, `worktree`, and `binary` events and
accumulates rows in React state. Before render, **repositoryScansLayout** applies
frontend-only **filter** rules (hide items under 1 MiB / 10 MiB when checkboxes
are unchecked) then **sort** rules (repos and children by size DESC). Checkbox
state (`showUnder1M`, `showUnder10M`) is independent per section. Re-sort runs
on every state update so newly arrived larger items move above smaller ones.

## Test Tree

```
repository-scans-layout/
├── sort/
│   ├── worktree-repos-desc/
│   ├── binary-repos-desc/
│   └── children-desc/
├── filter/
│   ├── binaries-hide-under-1m/
│   ├── binaries-show-when-checked/
│   ├── worktrees-hide-under-10m/
│   └── worktrees-show-when-checked/
└── resort/
    └── on-new-larger-item/
```

## Test Index

| Leaf | Op | Fixture |
|------|-----|---------|
| sort/worktree-repos-desc | sort-worktree-repos | testdata/fixture.json |
| sort/binary-repos-desc | sort-binary-repos | testdata/fixture.json |
| sort/children-desc | sort-linked-worktrees | testdata/fixture.json |
| filter/binaries-hide-under-1m | filter-binary-repos | testdata/fixture.json |
| filter/binaries-show-when-checked | filter-binary-repos | testdata/fixture.json |
| filter/worktrees-hide-under-10m | filter-worktree-repos | testdata/fixture.json |
| filter/worktrees-show-when-checked | filter-worktree-repos | testdata/fixture.json |
| resort/on-new-larger-item | resort-worktree-repos | testdata/fixture.json |

## How to Run

```sh
doctest vet ./tests/repository-scans-layout
doctest test ./tests/repository-scans-layout
```

Requires **Node.js** and network for first `npx tsx` fetch. The harness imports
`disk-usage-analyser-react/src/repositoryScansLayout.ts` (RED until implemented).

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

	harness := filepath.Join(projectRoot, "tests/repository-scans-layout/layout-harness.ts")
	cmd := exec.Command("npx", "--yes", "tsx", harness, "--op", req.Op, "--fixture", fixturePath)
	cmd.Dir = projectRoot
	out, err := cmd.CombinedOutput()
	output := string(out)
	if err != nil {
		return &Response{Output: output}, fmt.Errorf("layout harness: %w\n%s", err, output)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(out, &parsed); err != nil {
		return &Response{Output: output}, fmt.Errorf("parse harness JSON: %w\n%s", err, output)
	}

	return &Response{Output: output, JSON: parsed}, nil
}
```