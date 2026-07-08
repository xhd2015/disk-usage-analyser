## Expected

- Non-zero exit code.
- A non-nil error mentions the invalid threshold value.

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
		t.Fatal("expected threshold parse error, got nil")
	}
	msg := resp.Err.Error()
	if !strings.Contains(strings.ToLower(msg), "threshold") && !strings.Contains(strings.ToLower(msg), "foo") {
		t.Fatalf("error should mention invalid threshold; got %q", msg)
	}
}
```