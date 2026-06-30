# Scenario

**Decision**: worktree checkout sizing

```
sizing -> non-empty checkout reports size > 0
```

## Preconditions

- Sizing walks regular files under each worktree checkout path.
- Empty repos still have commit metadata but may have small non-zero size from `.git` objects in main checkout.

## Steps

1. Place known file content inside a worktree checkout.
2. Run `worktrees-scan` and inspect `size` fields.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "worktrees-scan"
	return nil
}
```