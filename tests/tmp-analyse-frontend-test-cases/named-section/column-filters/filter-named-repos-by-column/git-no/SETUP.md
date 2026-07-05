# Scenario

**Leaf**: Git filter `no` keeps only hits where `gitTracked !== true`

## Steps

1. Load fixture with tracked, untracked, and missing `gitTracked` hits.
2. Run `filter-named-column-filters` with `git: no`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```