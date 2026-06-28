# Scenario

**Feature**: build breakdown progress payload

```
# BuildBreakdownProgressPayload -> card accumulated size + active row fields
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related build breakdown progress payload leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	req.Reclaimable = true
	return nil
}
```
