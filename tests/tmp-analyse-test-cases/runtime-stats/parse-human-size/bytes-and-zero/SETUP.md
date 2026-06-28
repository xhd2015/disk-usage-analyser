# Scenario

**Feature**: parse "500 B" and "0 B"

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- "500 B" -> 500; "0 B" -> 0.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
import (
	"disk-usage-analyser/server"
)
func Setup(t *testing.T, req *Request) error {
	req.Op = "parse-human-size-bytes"
	req.HumanSize = "500 B"
	return nil
}
```
