# Scenario

**Feature**: card total monotonically increases across progress steps

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Simulated Go scan: card total never decreases between progress events.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "breakdown-sim-never-drops"
	req.Label = "Go"
	req.Reclaimable = true
	req.SimSequence = []breakdownSimStep{
		{ActiveIndex: 0, ActiveSize: 1000000000, ActiveCount: 100},
		{CompletedSizes: []int64{1000000000}, ActiveIndex: 1, ActiveSize: 500000000, ActiveCount: 50},
		{CompletedSizes: []int64{1000000000, 500000000}, ActiveIndex: 1, ActiveSize: 800000000, ActiveCount: 80},
	}
	return nil
}
```
