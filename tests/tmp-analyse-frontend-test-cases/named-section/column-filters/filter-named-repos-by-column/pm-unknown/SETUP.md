# Scenario

**Leaf**: PM filter `unknown` keeps hits with missing or empty `packageManager`

## Steps

1. Load fixture with npm, missing PM, and empty-string PM hits.
2. Run `filter-named-column-filters` with `pm: unknown`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```