# Scenario

**Feature**: enriched named SSE JSON payload fields on pass 2

```
named_enriched event JSON -> packageManager + pnpm/bun/shared size fields present on every enriched hit
```

## Preconditions

- Every `named_enriched` event must include all JSON keys even when values are zero/empty.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "named-scan"
	req.Name = "node_modules"
	return nil
}
```