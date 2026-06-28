# Scenario

**Feature**: npm progress includes breakdownPath for dynamic subdir

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- npm progress event includes breakdownPath for dynamic subdir row.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "breakdown-npm-path"
	req.BreakdownPath = "~/.npm/_cacache"
	return nil
}
```
