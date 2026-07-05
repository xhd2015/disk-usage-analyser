# Scenario

**Leaf**: truncation prefers breaking at a slash boundary before the suffix

## Steps

1. Load fixture where a slash-aligned cut is possible within the budget.
2. Run `truncate-path` op.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "truncate-path"
	return nil
}
```