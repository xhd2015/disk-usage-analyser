# Scenario

**Leaf**: leaf checkbox toggles update `binary-selected-total` size bar

## Steps

1. Set req.ScriptFile to binaries-select-total.js.
2. Scan binaries, toggle a leaf checkbox, verify selected total updates.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "binaries-select-total.js"
	return nil
}
```