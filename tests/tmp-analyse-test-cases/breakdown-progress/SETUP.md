# Scenario

**Feature**: live breakdown progress grouping

```
# ScanWithProgress -> BuildBreakdownProgressPayload -> progress SSE with breakdownIndex
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related live breakdown progress grouping leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Label == "" {
		req.Label = "Go"
	}
	if !req.Reclaimable {
		req.Reclaimable = true
	}
	if req.AccumulatedSize == 0 {
		req.AccumulatedSize = 0
	}
	if req.AccumulatedReclaimable == 0 {
		req.AccumulatedReclaimable = 0
	}
	return nil
}
```
