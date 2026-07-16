## Expected

- `ServerOptions.Dev` is false.
- `ServerOptions.DevIdleLife` is zero (flag ignored when not in dev mode).

## Side Effects

- Fake `StartServer` is invoked (normal server path without `--dev`).

## Errors

- No harness or CLI error.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected CLI error: %v", resp.Err)
	}
	if !resp.ServerWasStarted {
		t.Fatal("expected StartServer to be called")
	}
	if resp.ServerDev {
		t.Fatal("ServerOptions.Dev = true, want false")
	}
	if resp.DevIdleLife != 0 {
		t.Fatalf("DevIdleLife = %v, want 0 when --dev is absent", resp.DevIdleLife)
	}
}
```