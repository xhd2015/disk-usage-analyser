## Expected

- `GET /ping` returns HTTP 200 with body `pong`.
- After advancing the fake clock 5s with `DevIdleLife=0`, the port is still listening.

## Side Effects

- No idle shutdown log on stderr.

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
	if !resp.PortListening {
		t.Fatal("port closed with DevIdleLife=0, want still listening")
	}
}
```