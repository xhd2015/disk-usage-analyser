# Scenario

**Leaf**: SSE worktree events never carry `isMain=true`; main checkout sizing is only on repo events

## Preconditions

- Main repo has a linked worktree so at least one worktree event is emitted.

## Steps

1. Init `~/Projects/omit-main` and commit.
2. Add linked worktree `~/Projects/omit-main-wt` on branch `feature`.
3. Run `worktrees-scan`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/omit-main")
	gitInitialCommit(t, mainDir)
	wtDir := filepath.Join(req.HomeDir, "Projects", "omit-main-wt")
	gitWorktreeAdd(t, mainDir, wtDir, "feature")
	req.Op = "worktrees-scan"
	return nil
}
```