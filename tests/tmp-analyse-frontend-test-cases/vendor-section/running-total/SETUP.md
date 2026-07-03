# Scenario

**Leaf**: vendor title shows running total after scan via SSE hit events

## Steps

1. Set req.ScriptFile to vendor-running-total.js.
2. Click vendor scan, wait for done badge, verify running-total element shows formatted size.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "vendor-running-total.js"
	return nil
}
```
