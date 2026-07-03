# Scenario

**Leaf**: binaries title shows running total after scan via SSE hit events

## Steps

1. Set req.ScriptFile to binaries-running-total.js.
2. Click binaries scan, wait for done badge, verify running-total element shows formatted size.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-running-total.js"
	return nil
}
```
