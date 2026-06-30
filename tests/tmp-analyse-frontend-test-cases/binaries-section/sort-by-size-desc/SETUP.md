# Scenario

**Leaf**: binary repos and rows sorted by size DESC in DOM

## Steps

1. Set req.ScriptFile to binaries-sort-by-size-desc.js.
2. After scan, verify repo totals and row sizes are monotonic non-increasing.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-sort-by-size-desc.js"
	return nil
}
```