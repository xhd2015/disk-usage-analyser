# Scenario

**Feature**: npm dynamic breakdown path

```
# npm subdir scan -> progress event includes breakdownPath
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related npm dynamic breakdown path leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	req.BreakdownPath = "~/.npm/_cacache"
	return nil
}
```
