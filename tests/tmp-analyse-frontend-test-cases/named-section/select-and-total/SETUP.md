# Scenario

**Leaf**: leaf checkbox toggles update `node-modules-selected-total` size bar

## Steps

1. Set req.ScriptFile to named-select-total.js.
2. Scan node_modules, toggle a leaf checkbox, verify selected total updates.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-select-total.js"
	return nil
}
```
