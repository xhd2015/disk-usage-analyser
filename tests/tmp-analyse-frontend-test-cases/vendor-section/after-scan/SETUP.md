# Scenario

**Leaf**: vendor scan populates repo-grouped tree independently

## Steps

1. Set req.ScriptFile to vendor-after-scan.js.
2. Run vendor scan and verify tree renders with size, name, path, repo columns.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "vendor-after-scan.js"
	return nil
}
```
