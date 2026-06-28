# Scenario

**Feature**: completed row 0 frozen while row 1 active

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Row 0 frozen at 3.1GB; row 1 partial 0.8GB; card total 3.9GB.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "breakdown-payload-second-row"
	req.Label = "Go"
	req.CompletedSizes = []int64{3100000000}
	req.CompletedCounts = []int64{2000}
	req.ActiveIndex = 1
	req.ActiveSize = 800000000
	req.ActiveCount = 1500
	return nil
}
```
