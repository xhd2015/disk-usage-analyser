# Scenario

**Leaf**: worktree repos sorted by main checkout size DESC

## Steps

1. Load fixture with three repos (5 MB, 20 MB, 10 MB).
2. Run `sort-worktree-repos`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "sort-worktree-repos"
	return nil
}
```