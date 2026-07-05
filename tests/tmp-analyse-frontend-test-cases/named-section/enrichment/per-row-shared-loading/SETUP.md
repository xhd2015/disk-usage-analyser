# Scenario

**Leaf**: per-row shared column shows loading while enrichment is pending

## Steps

1. Set req.ScriptFile to per-row-shared-loading.js.
2. Run node_modules scan; during enrichment, un-enriched rows expose `node-modules-shared-loading-*`.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "per-row-shared-loading.js"
	return nil
}
```