# Scenario

**Leaf**: worktrees title shows running total after scan via SSE hit events

## Steps

1. Set req.ScriptFile to worktrees-running-total.js.
2. Click worktrees scan, wait for done badge, verify running-total element shows formatted size.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-running-total.js"
	return nil
}
```
