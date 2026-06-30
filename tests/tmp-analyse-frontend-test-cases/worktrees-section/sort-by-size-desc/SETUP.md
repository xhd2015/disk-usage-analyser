# Scenario

**Leaf**: worktree repos and linked children sorted by size DESC in DOM

## Steps

1. Set req.ScriptFile to worktrees-sort-by-size-desc.js.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-sort-by-size-desc.js"
	return nil
}
```