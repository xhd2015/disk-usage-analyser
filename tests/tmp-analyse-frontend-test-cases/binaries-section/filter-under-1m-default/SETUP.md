# Scenario

**Leaf**: `<1M` checkbox unchecked by default; binaries under 1 MiB hidden

## Steps

1. Set req.ScriptFile to binaries-filter-under-1m-default.js.
2. Run binaries scan and verify filter checkbox + visible row sizes.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-filter-under-1m-default.js"
	return nil
}
```