# Scenario

**Feature**: completed breakdown row sizes retained in payload

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Completed row 0 size retained; active row 1 adds partial.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "breakdown-sim-rows-retained"
	req.Label = "Go"
	req.Reclaimable = true
	req.SimSequence = []breakdownSimStep{
		{CompletedSizes: []int64{3100000000}, ActiveIndex: 1, ActiveSize: 200000000, ActiveCount: 10},
	}
	return nil
}
```
