# Scenario

**Feature**: collect podman stats via mock runner

```
# tmp-analyse runtime/breakdown pipeline under test
HandleTmpAnalyse -> SSE events -> frontend cards
```

## Preconditions
- Mock podman CLI returns filtered runtime items.

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
	req.Op = "collect-runtime-podman"
	req.Runtime = "podman"
	req.FixtureFile = "testdata/podman-system-df.ndjson"
	data, _ := os.ReadFile(req.FixtureFile)
	req.MockOutput = string(data)
	return nil
}
```
