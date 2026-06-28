# Scenario

**Feature**: first breakdown row active during scan

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- First row scanning: card size equals active partial only.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "breakdown-payload-first-row"
	req.Label = "Go"
	req.ActiveIndex = 0
	req.ActiveSize = 500000000
	req.ActiveCount = 100
	return nil
}
```
