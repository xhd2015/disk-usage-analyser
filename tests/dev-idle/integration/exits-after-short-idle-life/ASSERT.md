## Expected

- `GET /ping` returns HTTP 200 with body `pong`.
- After advancing past `DevIdleLife=2s`, the listening port is closed (TCP dial fails).

## Side Effects

- `server.ListenAndServe` returns after idle shutdown (not a hang).

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
		t.Fatal("port still listening after idle shutdown, want closed")
	}
}
```