# Scenario

**Leaf**: linked worktrees within a repo sorted by hit size DESC

## Steps

1. Load fixture with linked worktrees 2 MB, 8 MB, 5 MB.
2. Run `sort-linked-worktrees`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "sort-linked-worktrees"
	return nil
}
```