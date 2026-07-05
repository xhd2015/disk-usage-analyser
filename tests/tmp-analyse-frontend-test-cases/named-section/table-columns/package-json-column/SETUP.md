# Scenario

**Feature**: package.json read-only checkbox column before PM in node_modules table

## Steps

1. Set req.ScriptFile to table-columns-package-json-column.js.
2. Run node_modules scan and verify header + checkbox cells.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "table-columns-package-json-column.js"
	return nil
}
```