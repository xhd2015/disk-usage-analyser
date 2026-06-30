# Scenario

**Leaf**: checking `<10M` shows all worktree sizes

## Steps

1. Same fixture as hide-under-10m with `showUnder10M: true`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-worktree-repos"
	return nil
}
```