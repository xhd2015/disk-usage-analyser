# Scenario

**Leaf**: pipelined enrichment emits first `named_enriched` before `scan_complete`

## Steps

1. Create two git repos each with `node_modules/` (multi-hit fixture).
2. Run `named-scan`.

```go
func Setup(t *testing.T, req *Request) error {
	nodeModulesRepo(t, req.HomeDir, "Projects/pipeline-app-a")
	nodeModulesRepo(t, req.HomeDir, "Projects/pipeline-app-b")
	return nil
}
```