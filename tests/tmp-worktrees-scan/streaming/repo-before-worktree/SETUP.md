# Scenario

**Leaf**: `repo` SSE event arrives before the first linked `worktree` event for that scan

## Preconditions

- Fixture home contains a git repo with a linked worktree.

## Steps

1. Create `~/Projects/order-wt` with initial commit.
2. Add linked worktree `~/Projects/order-wt-linked` on branch `feature`.
3. Run `worktrees-sse-order`.

## Context

- The frontend uses `repo` events to create parent tree nodes before child worktree rows stream in.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/order-wt")
	gitInitialCommit(t, mainDir)
	wtDir := filepath.Join(req.HomeDir, "Projects", "order-wt-linked")
	gitWorktreeAdd(t, mainDir, wtDir, "feature")
	req.Op = "worktrees-sse-order"
	return nil
}
```