## Expected

- Harness starts `--dev` server and discovers listening port.
- After explicit teardown (no playwright), TCP dial to that port fails.
- Listener PID resolved via `lsof` while server was up is no longer alive.

## Side Effects

- None beyond normal harness server start/stop.

## Errors

- Port still accepting connections after teardown.
- Listener PID still responds to signal 0 after teardown.

## Exit Code

- 0 when port closed and listener dead; non-zero on orphan.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.ServerPort == "" {
		t.Fatal("expected ServerPort to be set")
	}
	if resp.ListenerPID <= 0 {
		t.Fatalf("expected listener PID on port %s before teardown, got %d", resp.ServerPort, resp.ListenerPID)
	}
	if resp.PortListeningAfterTeardown {
		t.Fatalf("port %s still listening after teardown (orphaned --dev server?)", resp.ServerPort)
	}
	if resp.ListenerAliveAfterTeardown {
		t.Fatalf("listener pid %d still alive after teardown (expected compiled binary killed, not go wrapper only)", resp.ListenerPID)
	}
}
```