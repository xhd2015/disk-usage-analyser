# Scenario

**Feature**: runtime stats grouping

```
# CollectRuntimeStats -> ParseSystemDFJSON -> FilterRuntimeItems -> runtimeItems on location event
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related runtime stats grouping leaves.

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
