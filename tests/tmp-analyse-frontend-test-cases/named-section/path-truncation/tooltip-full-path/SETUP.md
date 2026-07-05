# Scenario

**Leaf**: path cell ellipsis styling, full path in DOM, tooltip shows full path on hover

## Steps

1. Set req.ScriptFile to path-truncation-tooltip.js.
2. Run node_modules scan and inspect first path cell.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "path-truncation-tooltip.js"
	return nil
}
```