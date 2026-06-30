# Scenario

**Leaf**: default filter hides worktrees under 10 MiB

## Steps

1. Load fixture: main 5 MB repo hidden; main 20 MB repo visible with only 12 MB linked child.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-worktree-repos"
	return nil
}
```