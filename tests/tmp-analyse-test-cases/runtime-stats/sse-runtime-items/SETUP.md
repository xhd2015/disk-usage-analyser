# Scenario

**Feature**: SSE location carries runtimeItems

```
# HandleTmpAnalyse -> CollectRuntimeStats -> location event runtimeItems
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related SSE location carries runtimeItems leaves.

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
