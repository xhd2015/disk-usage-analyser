# Scenario

**Leaf**: checking `<1M` reveals binaries under 1 MiB when they exist

## Steps

1. Set req.ScriptFile to binaries-filter-show-under-1m.js.
2. Compare row counts before and after checking `binary-show-under-1m`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-filter-show-under-1m.js"
	return nil
}
```