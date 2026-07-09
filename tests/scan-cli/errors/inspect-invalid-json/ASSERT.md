## Expected

- Non-zero exit code.
- Error mentions JSON decode / invalid payload.

## Exit Code

- Non-zero

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0 (err=%v)", resp.Err)
	}
	if resp.Err == nil {
		t.Fatal("expected invalid JSON error, got nil")
	}
	msg := strings.ToLower(resp.Err.Error())
	if !strings.Contains(msg, "json") && !strings.Contains(msg, "decode") && !strings.Contains(msg, "invalid") {
		t.Fatalf("error should mention JSON/decode failure; got %q", msg)
	}
}
```
