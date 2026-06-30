# Scenario

**Leaf**: worktree hits use `~/...` paths for checkouts under the fixture home

## Preconditions

- Repo lives at `~/Projects/tilde-app`.

## Steps

1. Init `~/Projects/tilde-app` and commit.
2. Run `worktrees-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	if !gitAvailable(t) {
		return nil
	}
	mainDir := repoUnderHome(t, req.HomeDir, "Projects/tilde-app")
	gitInitialCommit(t, mainDir)
	req.Op = "worktrees-scan"
	return nil
}
```