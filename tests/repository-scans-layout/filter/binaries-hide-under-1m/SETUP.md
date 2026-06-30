# Scenario

**Leaf**: default filter hides binaries under 1 MiB and tiny repos

## Steps

1. Load fixture with 500 KB-only repo and mixed repo (800 KB + 2 MB).
2. Run `filter-binary-repos` with `showUnder1M: false`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-binary-repos"
	return nil
}
```