# Scenario

**Leaf**: combined Git `yes` and PM `npm` filters intersect

## Steps

1. Load fixture with tracked/untracked hits across npm and pnpm.
2. Run `filter-named-column-filters` with `git: yes` and `pm: npm`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-named-column-filters"
	return nil
}
```