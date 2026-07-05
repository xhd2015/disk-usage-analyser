# Scenario

**Leaf**: after node_modules scan, table grid spans the full card body width

## Steps

1. Set req.ScriptFile to table-columns-full-width.js.
2. Run node_modules scan at 1280px viewport.
3. Compare column header width to card body width.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "table-columns-full-width.js"
	return nil
}
```