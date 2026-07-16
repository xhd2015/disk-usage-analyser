## Expected

- SSE stream delivers at least 6 `tick` events before `done`.
- After stream ends and fake clock advances 3s without traffic, the listening port is closed.

## Side Effects

- `server.ListenAndServe` returns after post-stream idle shutdown.

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
	if !resp.PortClosed {
		t.Fatal("port still listening after post-stream idle, want closed")
	}
}
```