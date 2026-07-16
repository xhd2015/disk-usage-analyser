## Expected

- CLI returns a non-nil error for the invalid duration string.
- Fake `StartServer` is **not** invoked.

## Errors

- `resp.Err` must be non-nil (invalid duration rejected).

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err == nil {
		t.Fatal("expected error for invalid --dev-idle-life value, got nil")
	}
	if resp.ServerWasStarted {
		t.Fatal("StartServer must not be called when duration parse fails")
	}
}
```