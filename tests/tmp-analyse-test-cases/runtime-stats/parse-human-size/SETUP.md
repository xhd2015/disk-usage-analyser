# Scenario

**Feature**: parse human size strings

```
# ParseSystemDFJSON -> ParseHumanSize -> byte counts
HandleTmpAnalyse -> SSE progress/location events -> frontend breakdown rows
```

## Preconditions
- Grouping node for related parse human size strings leaves.

## Steps
1. Descendant leaf sets `req.Op` and scenario-specific fields.

## Context
- See leaf ASSERT.md for concrete expectations.

```go
func Setup(t *testing.T, req *Request) error {
	if req.HumanSize == "" {
		req.HumanSize = "1 MB"
	}
	return nil
}
```
