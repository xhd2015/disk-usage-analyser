# Scenario

**Leaf**: worktrees tree content is flush left within the card

## Steps

1. Set req.ScriptFile to worktrees-left-aligned.js.
2. Run worktrees scan, wait for done, check computed text-align on worktrees-tree.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-left-aligned.js"
	return nil
}
```