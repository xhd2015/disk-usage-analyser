# Scenario

**Leaf**: node_modules results sorted largest first

## Steps

1. Set req.ScriptFile to named-sort-by-size-desc.js.
2. Scan, verify repo group totals and within-repo row sizes are monotonic non-increasing.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-sort-by-size-desc.js"
	return nil
}
```
