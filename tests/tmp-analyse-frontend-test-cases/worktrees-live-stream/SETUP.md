# Scenario

**Leaf**: worktree row count grows during an in-progress scan

## Steps

1. Set req.ScriptFile to worktrees-live-stream.js.
2. Start worktrees scan and poll row count while scanning badge is visible.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-live-stream.js"
	return nil
}
```