## Expected

- SSE stream delivers at least 6 `tick` events over 3s (500ms interval).
- After the stream completes, the listening port is still open (`DevIdleLife=2s` exceeded only
  because each `sendSSEEvent` touch kept the watch alive).

## Side Effects

- No idle shutdown while SSE events are flowing.

## Errors

- No harness error.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.SSEEventCount < req.StreamTicks {
		t.Fatalf("sse tick count = %d, want >= %d", resp.SSEEventCount, req.StreamTicks)
	}
	if !resp.PortListening {
		t.Fatal("port closed during active SSE stream, want still listening")
	}
}
```