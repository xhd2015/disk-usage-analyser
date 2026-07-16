# Scenario

**Bug**: no orphaned `--dev` listener after harness teardown

```
# start server, skip playwright, assert port closed and listener PID dead
doctest Run -> startDevServer -> teardown -> port closed, listener PID not alive
```

## Preconditions

- Server starts and prints `Serving directory preview at http://localhost:<port>`.
- `lsof` is available to resolve listener PID while port is open.

## Steps

1. Leaf inherits `req.TeardownOnly = true` from parent grouping.
2. `Run` starts server, runs teardown, records post-teardown listener state.

## Context

- RED until `Run` uses compiled binary and kills listener/process group.

```go
import (
	"os/exec"
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	_, err := exec.LookPath("lsof")
	if err != nil {
		t.Skip("lsof not available for listener PID resolution")
	}
	req.TeardownOnly = true
	return nil
}
```