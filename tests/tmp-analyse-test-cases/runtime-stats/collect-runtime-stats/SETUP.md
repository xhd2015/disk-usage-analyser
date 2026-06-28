# Scenario

**Feature**: collect runtime via CLI

```
# CollectRuntimeStats(runtime) -> mock command runner -> filtered items
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related collect runtime via CLI leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Runtime == "" {
		req.Runtime = "docker"
	}
	return nil
}
```
