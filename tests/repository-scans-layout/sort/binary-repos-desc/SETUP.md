# Scenario

**Leaf**: binary repos sorted by sum of visible binary sizes DESC

## Steps

1. Load fixture with repo totals 3 MB and 15 MB.
2. Run `sort-binary-repos`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "sort-binary-repos"
	return nil
}
```