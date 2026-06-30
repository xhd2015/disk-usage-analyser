# Scenario

**Leaf**: `<10M` checkbox unchecked by default; worktrees under 10 MiB hidden

## Steps

1. Set req.ScriptFile to worktrees-filter-under-10m-default.js.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-filter-under-10m-default.js"
	return nil
}
```