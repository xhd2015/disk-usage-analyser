# Scenario

**Leaf**: long path gets ellipsis prefix while suffix stays fully visible

## Steps

1. Load fixture with a long path and small maxVisibleChars.
2. Run `truncate-path` op.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "truncate-path"
	return nil
}
```