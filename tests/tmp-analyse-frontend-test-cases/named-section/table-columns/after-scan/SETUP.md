# Scenario

**Leaf**: after node_modules scan, table header and pkgmgr/shared columns render

## Steps

1. Set req.ScriptFile to table-columns-after-scan.js.
2. Run node_modules scan and verify column cells.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "table-columns-after-scan.js"
	return nil
}
```