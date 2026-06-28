# Scenario

**Feature**: parse compact "1.2GB"

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- "1.2GB" parses to 1200000000 bytes.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-human-size-gb"
	req.HumanSize = "1.2GB"
	return nil
}
```
