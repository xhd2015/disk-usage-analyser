# Scenario

**Leaf**: two main-only repositories each produce a repo event with main checkout size

## Preconditions

- Each repo has an initial commit.

## Steps

1. Init `~/Projects/alpha` and `~/Work/beta` with commits.
2. Run `worktrees-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	alpha := repoUnderHome(t, req.HomeDir, "Projects/alpha")
	gitInitialCommit(t, alpha)
	beta := repoUnderHome(t, req.HomeDir, "Work/beta")
	gitInitialCommit(t, beta)
	req.Op = "worktrees-scan"
	return nil
}
```