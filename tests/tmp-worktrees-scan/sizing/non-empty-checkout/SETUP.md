# Scenario

**Leaf**: worktree hit for a checkout with files reports `size > 0` and `fileCount > 0`

## Preconditions

- Linked worktree contains additional tracked files beyond the initial commit.

## Steps

1. Init `~/Projects/size-app` and commit.
2. Add linked worktree `~/Projects/size-wt` on branch `feature`.
3. Write `payload.bin` (4096 bytes) into the linked worktree.
4. Run `worktrees-scan`.

```go
import (
	"path/filepath"
)

func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/size-app")
	gitInitialCommit(t, mainDir)
	wtDir := filepath.Join(req.HomeDir, "Projects", "size-wt")
	gitWorktreeAdd(t, mainDir, wtDir, "feature")
	writeFile(t, wtDir, "payload.bin", make([]byte, 4096))
	req.Op = "worktrees-scan"
	return nil
}
```