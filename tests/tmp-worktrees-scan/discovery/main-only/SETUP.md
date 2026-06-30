# Scenario

**Leaf**: single main repo emits repo event with main checkout size and zero worktree events

## Preconditions

- Repo has an initial commit so `git worktree list` succeeds.

## Steps

1. Init `~/Projects/solo` and commit.
2. Run `worktrees-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/solo")
	gitInitialCommit(t, mainDir)
	req.Op = "worktrees-scan"
	return nil
}
```