# Scenario

**Feature**: reject invalid human size

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Invalid size string returns error.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-human-size-invalid"
	req.HumanSize = "not-a-size"
	return nil
}
```
