# Scenario

**Leaf**: Git filter `yes` keeps only `gitTracked === true` hits

## Steps

1. Load fixture with tracked, untracked, and missing `gitTracked` hits.
2. Run `filter-named-column-filters` with `git: yes`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```