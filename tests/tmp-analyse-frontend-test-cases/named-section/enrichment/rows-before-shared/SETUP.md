# Scenario

**Leaf**: node_modules rows visible at scan_complete; shared column populated or stays 0 B after enrichment

## Steps

1. Set req.ScriptFile to rows-before-shared.js.
2. Run node_modules scan and poll row/done/shared timing.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "rows-before-shared.js"
	return nil
}
```