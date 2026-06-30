# Scenario

**Leaf**: worktrees scan populates repo tree with worktree rows and sizes

## Steps

1. Set req.ScriptFile to worktrees-after-scan.js.
2. Click worktrees scan, wait for done, verify tree rows.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-after-scan.js"
	return nil
}
```