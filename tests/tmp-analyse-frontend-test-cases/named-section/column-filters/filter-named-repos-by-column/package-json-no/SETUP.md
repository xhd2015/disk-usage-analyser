# Scenario

**Leaf**: package.json filter `no` keeps only hits where `hasPackageJson !== true`

## Steps

1. Load fixture with true, false, and missing `hasPackageJson` hits.
2. Run `filter-named-column-filters` with `packageJson: no`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```