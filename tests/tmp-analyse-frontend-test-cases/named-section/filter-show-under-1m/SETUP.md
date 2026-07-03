# Scenario

**Leaf**: toggle filter checkbox reveals small node_modules entries

## Steps

1. Set req.ScriptFile to named-filter-show-under-1m.js.
2. Scan, count rows before/after checking filter, verify count increases when small entries exist.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-filter-show-under-1m.js"
	return nil
}
```
