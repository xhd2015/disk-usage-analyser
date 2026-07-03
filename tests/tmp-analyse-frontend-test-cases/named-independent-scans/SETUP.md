# Scenario

**Leaf**: node_modules and vendor scans run concurrently without blocking each other

## Steps

1. Set req.ScriptFile to named-independent-scans.js.
2. Start node_modules scan, verify vendor scan button is still visible and can be started simultaneously.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-independent-scans.js"
	return nil
}
```
