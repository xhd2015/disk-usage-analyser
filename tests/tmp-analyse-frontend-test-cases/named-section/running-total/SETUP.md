# Scenario

**Leaf**: node_modules title shows running total after scan via SSE hit events

## Steps

1. Set req.ScriptFile to node-modules-running-total.js.
2. Click node_modules scan, wait for done badge, verify running-total element shows formatted size.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "node-modules-running-total.js"
	return nil
}
```
