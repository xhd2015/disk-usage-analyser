# Scenario

**Leaf**: main repo reports main checkout size on repo event plus linked worktree child only

## Preconditions

- Real `git` is available.
- Main repo has an initial commit before adding a linked worktree.

## Steps

1. Init `~/Projects/foo` and commit.
2. Add linked worktree `~/Projects/foo-wt` on branch `feature`.
3. Run `worktrees-scan`.

## Context

- Main checkout path is the repo root; linked path is the sibling worktree directory.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/foo")
	gitInitialCommit(t, mainDir)
	wtDir := filepath.Join(req.HomeDir, "Projects", "foo-wt")
	gitWorktreeAdd(t, mainDir, wtDir, "feature")
	req.Op = "worktrees-scan"
	return nil
}
```