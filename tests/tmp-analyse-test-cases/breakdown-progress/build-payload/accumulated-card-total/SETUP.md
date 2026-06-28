# Scenario

**Feature**: card size equals completed rows plus active partial

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- size=4.2GB (accumulated card), breakdownIndex=1, breakdownSize=1.1GB.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "breakdown-payload-accumulated"
	req.Label = "Go"
	req.CompletedSizes = []int64{3100000000}
	req.CompletedCounts = []int64{1000}
	req.ActiveIndex = 1
	req.ActiveSize = 1100000000
	req.ActiveCount = 500
	req.AccumulatedSize = 5000000000
	req.AccumulatedReclaimable = 5000000000
	req.Reclaimable = true
	return nil
}
```
