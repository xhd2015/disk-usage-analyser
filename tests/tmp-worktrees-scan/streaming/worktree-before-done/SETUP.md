# Scenario

**Leaf**: at least one linked `worktree` SSE event arrives before `done`

## Preconditions

- Fixture home contains a git repo with a linked worktree (main-only repos emit zero worktree events).

## Steps

1. Create `~/Projects/stream-wt` with initial commit.
2. Add linked worktree `~/Projects/stream-wt-linked` on branch `feature`.
3. Run `worktrees-sse-order` against `/api/tmp-worktrees-scan`.

## Context

- Buffered scans that emit only `done` fail this leaf.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/stream-wt")
	gitInitialCommit(t, mainDir)
	wtDir := filepath.Join(req.HomeDir, "Projects", "stream-wt-linked")
	gitWorktreeAdd(t, mainDir, wtDir, "feature")
	req.Op = "worktrees-sse-order"
	return nil
}
```