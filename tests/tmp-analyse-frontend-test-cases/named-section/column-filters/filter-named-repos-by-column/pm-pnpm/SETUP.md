# Scenario

**Leaf**: PM filter `pnpm` keeps only pnpm rows

## Steps

1. Load fixture with npm, pnpm, and yarn hits.
2. Run `filter-named-column-filters` with `pm: pnpm`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```