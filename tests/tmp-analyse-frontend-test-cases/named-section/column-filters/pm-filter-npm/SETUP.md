# Scenario

**Leaf**: PM filter `npm` shows only npm package-manager rows when mixed PMs exist

## Steps

1. Set req.ScriptFile to column-filters-pm-filter-npm.js.
2. Scan, detect mixed PMs, select PM=npm, verify visible PM cells are all npm.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "column-filters-pm-filter-npm.js"
	return nil
}
```