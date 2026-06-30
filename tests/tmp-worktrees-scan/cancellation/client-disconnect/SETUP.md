# Scenario

**Leaf**: closing the SSE connection after the first repo event aborts the scan

## Preconditions

- Fixture contains at least one repo slow enough to stream multiple worktrees, or multiple repos.

## Steps

1. Init `~/Projects/disc-a` and `~/Projects/disc-b` with commits.
2. Run `worktrees-disconnect`.

```go
func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	a := repoUnderHome(t, req.HomeDir, "Projects/disc-a")
	gitInitialCommit(t, a)
	b := repoUnderHome(t, req.HomeDir, "Projects/disc-b")
	gitInitialCommit(t, b)
	req.Op = "worktrees-disconnect"
	return nil
}
```