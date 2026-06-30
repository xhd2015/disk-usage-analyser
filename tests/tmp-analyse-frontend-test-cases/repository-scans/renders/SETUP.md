# Scenario

**Leaf**: Repository Scans group and both subsections render with scan buttons

## Steps

1. Set req.ScriptFile to repository-scans-renders.js.
2. Script checks section heading, worktrees section, binaries section, and scan buttons.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "repository-scans-renders.js"
	return nil
}
```