# Scenario

**Leaf**: small node_modules entries hidden by default

## Steps

1. Set req.ScriptFile to named-filter-under-1m-default.js.
2. Scan node_modules, verify filter checkbox is unchecked and sub-1M entries hidden.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-filter-under-1m-default.js"
	return nil
}
```
