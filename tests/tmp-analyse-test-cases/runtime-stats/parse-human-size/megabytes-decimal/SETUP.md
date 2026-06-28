# Scenario

**Feature**: parse "16.43 MB" to bytes

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- "16.43 MB" parses to 16430000 bytes (decimal MB).

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-human-size-mb"
	req.HumanSize = "16.43 MB"
	return nil
}
```
