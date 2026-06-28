# Scenario

**Feature**: missing binary returns empty slice without error

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Missing CLI binary: empty runtimeItems, nil error.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-runtime-missing"
	req.Runtime = "docker"
	req.MockNotFound = true
	return nil
}
```
