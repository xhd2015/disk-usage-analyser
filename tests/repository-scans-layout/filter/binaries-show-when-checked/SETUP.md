# Scenario

**Leaf**: checking `<1M` shows all binary sizes

## Steps

1. Same fixture as hide-under-1m with `showUnder1M: true`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "filter-binary-repos"
	return nil
}
```