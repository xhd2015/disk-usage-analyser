---
label: slow
explanation: real wall-clock idle wait (~12s; 10s watchdog tick) for subprocess stderr capture
---

## Expected Output

```text
[dev] no requests for 1s; shutting down
```

## Expected

- After idle expiry, captured stderr contains the shutdown log line with the configured idle duration.
- Subprocess exits and port is closed.

## Side Effects

- Shutdown log written to stderr before process exit.

## Errors

- No harness error.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if !resp.PortClosed {
		t.Fatalf("port %d still listening after idle shutdown, want closed", resp.Port)
	}
	if !resp.ProcessExited {
		t.Fatal("subprocess still running after idle shutdown, want exited")
	}
	if !strings.Contains(resp.Stderr, "[dev] no requests for") {
		t.Fatalf("stderr missing idle shutdown prefix; got:\n%s", resp.Stderr)
	}
	if !strings.Contains(resp.Stderr, "; shutting down") {
		t.Fatalf("stderr missing shutdown suffix; got:\n%s", resp.Stderr)
	}
	assert.Output(t, extractShutdownLine(resp.Stderr), `---
version: 2
---
[dev] no requests for 1s; shutting down`)
}

func extractShutdownLine(stderr string) string {
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[dev] no requests for") {
			return line
		}
	}
	return stderr
}
```