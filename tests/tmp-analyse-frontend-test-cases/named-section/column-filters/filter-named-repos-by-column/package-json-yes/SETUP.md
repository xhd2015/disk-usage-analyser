# Scenario

**Leaf**: package.json filter `yes` keeps only `hasPackageJson === true` hits

## Steps

1. Load fixture with true, false, and missing `hasPackageJson` hits.
2. Run `filter-named-column-filters` with `packageJson: yes`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```