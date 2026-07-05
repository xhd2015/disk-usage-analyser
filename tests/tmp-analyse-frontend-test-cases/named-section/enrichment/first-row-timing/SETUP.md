# Scenario

**Leaf**: first node_modules row appears within 10s of scan click (pass 1 responsiveness)

## Steps

1. Set req.ScriptFile to first-row-timing.js.
2. Click scan; assert first UI row within 10s budget.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "first-row-timing.js"
	return nil
}
```