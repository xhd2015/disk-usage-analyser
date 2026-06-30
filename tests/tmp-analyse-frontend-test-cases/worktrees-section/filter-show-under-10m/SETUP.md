# Scenario

**Leaf**: checking `<10M` reveals worktrees under 10 MiB when they exist

## Steps

1. Set req.ScriptFile to worktrees-filter-show-under-10m.js.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "worktrees-filter-show-under-10m.js"
	return nil
}
```