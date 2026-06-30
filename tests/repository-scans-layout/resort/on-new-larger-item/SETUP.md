# Scenario

**Leaf**: inserting a 25 MB repo moves it above existing 10 MB and 5 MB repos

## Steps

1. Sort initial repos (10 MB, 5 MB).
2. Insert 25 MB repo and re-sort.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "resort-worktree-repos"
	return nil
}
```