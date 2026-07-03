# Scenario

**Leaf**: node_modules scan populates repo-grouped tree with size, name, path, repo columns

## Steps

1. Set req.ScriptFile to named-after-scan.js.
2. Run node_modules scan and verify tree structure.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-after-scan.js"
	return nil
}
```
