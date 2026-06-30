# Scenario

**Leaf**: repository subsection scans do not block the page-level cache location scan

## Steps

1. Set req.ScriptFile to independent-scan-controls.js.
2. Start worktrees scan, then verify page-level Start Scan still works.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "independent-scan-controls.js"
	return nil
}
```