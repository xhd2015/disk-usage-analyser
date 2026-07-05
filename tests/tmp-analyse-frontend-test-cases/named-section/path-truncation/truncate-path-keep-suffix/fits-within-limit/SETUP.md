# Scenario

**Leaf**: short path fits within maxVisibleChars unchanged

## Steps

1. Load fixture with a path shorter than the limit.
2. Run `truncate-path` op.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "truncate-path"
	return nil
}
```