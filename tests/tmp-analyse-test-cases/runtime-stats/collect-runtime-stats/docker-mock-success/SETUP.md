# Scenario

**Feature**: collect docker stats via mock runner

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Mock docker CLI returns filtered Images + Build Cache.

## Steps
1. Set req.Op and scenario fields.
2. Run executes the targeted server function.

## Context
- Platform-independent pure function or mock-runner test.

```go
import (
	"os"
)
func Setup(t *testing.T, req *Request) error {
	req.Op = "collect-runtime-docker"
	req.Runtime = "docker"
	req.FixtureFile = "testdata/docker-system-df.ndjson"
	data, _ := os.ReadFile(req.FixtureFile)
	req.MockOutput = string(data)
	return nil
}
```
