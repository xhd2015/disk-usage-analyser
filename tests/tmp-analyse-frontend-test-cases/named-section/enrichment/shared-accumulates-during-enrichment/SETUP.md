# Scenario

**Leaf**: shared column total accumulates incrementally as named_enriched events arrive

## Steps

1. Set req.ScriptFile to shared-accumulates.js.
2. During enrichment, sum of per-row shared values increases at least twice before done.

```go
func Setup(t *testing.T, req *Request) error {
	req.ScriptFile = "shared-accumulates.js"
	return nil
}
```