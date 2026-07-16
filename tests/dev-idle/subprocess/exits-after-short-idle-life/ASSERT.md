---
label: slow
explanation: real wall-clock idle wait (~12s; 10s watchdog tick) for subprocess e2e
---

## Expected

- `GET /ping` returns HTTP 200 with body `pong`.
- After 3s real sleep past `DevIdleLife=2s`, TCP dial to the server port fails.
- The subprocess process has exited (not still running).

## Side Effects

- Subprocess terminates after idle shutdown (not a hang).

## Errors

- No harness error.

```go
import (
	"net/http"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.PingStatus != http.StatusOK {
		t.Fatalf("ping status = %d, want 200", resp.PingStatus)
	}
	if resp.PingBody != "pong" {
		t.Fatalf("ping body = %q, want %q", resp.PingBody, "pong")
	}
	if !resp.PortClosed {
		t.Fatalf("port %d still listening after idle shutdown, want closed", resp.Port)
	}
	if !resp.ProcessExited {
		t.Fatal("subprocess still running after idle shutdown, want exited")
	}
}
```