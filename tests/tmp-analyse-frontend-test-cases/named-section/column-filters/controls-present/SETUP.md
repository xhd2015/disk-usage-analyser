# Scenario

**Leaf**: column filter controls render with default `all` selection after scan

## Steps

1. Set req.ScriptFile to column-filters-controls-present.js.
2. Scan node_modules and verify three filter controls exist with default `all`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "column-filters-controls-present.js"
	return nil
}
```