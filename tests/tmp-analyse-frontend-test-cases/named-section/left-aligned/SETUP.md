# Scenario

**Leaf**: section heading and controls are left-aligned within the card

## Steps

1. Set req.ScriptFile to named-left-aligned.js.
2. Scan node_modules, verify tree text-align is left.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "named-left-aligned.js"
	return nil
}
```
