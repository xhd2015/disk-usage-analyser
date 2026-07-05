# Scenario

**Leaf**: default filters (`all` for each column) show every hit

## Steps

1. Load fixture with mixed Git, package.json, and PM values across two repos.
2. Run `filter-named-column-filters` with all filters set to `all`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```