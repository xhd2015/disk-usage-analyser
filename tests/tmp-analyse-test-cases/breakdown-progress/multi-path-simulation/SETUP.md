# Scenario

**Feature**: multi-path scan simulation

```
# completed breakdown rows + active row -> monotonic card total
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related multi-path scan simulation leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	req.Label = "Go"
	req.Reclaimable = true
	return nil
}
```
