## Expected

- `ServerOptions.Dev` is true.
- `ServerOptions.DevIdleLife` equals `10 * time.Minute`.

## Side Effects

- Fake `StartServer` is invoked (server path reached).

## Errors

- No harness or CLI error.

```go
import (
	"testing"
	"time"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected CLI error: %v", resp.Err)
	}
	if !resp.ServerWasStarted {
		t.Fatal("expected StartServer to be called for --dev")
	}
	if !resp.ServerDev {
		t.Fatal("ServerOptions.Dev = false, want true")
	}
	if resp.DevIdleLife != 10*time.Minute {
		t.Fatalf("DevIdleLife = %v, want 10m", resp.DevIdleLife)
	}
}
```