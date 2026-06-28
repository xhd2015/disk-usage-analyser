# Scenario

**Feature**: parse system df NDJSON

```
# docker system df --format json -> ParseSystemDFJSON -> TmpRuntimeItem slice
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related parse system df NDJSON leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	req.FixtureFile = "testdata/docker-system-df.ndjson"
	return nil
}
```
