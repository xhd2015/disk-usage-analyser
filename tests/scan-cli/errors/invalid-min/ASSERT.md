## Expected

- Non-zero exit code.
- A non-nil error mentions the invalid min value (`foo` and/or `min`).

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
		t.Fatal("expected min parse error, got nil")
	}
	msg := strings.ToLower(resp.Err.Error())
	if !strings.Contains(msg, "min") && !strings.Contains(msg, "foo") && !strings.Contains(msg, "size") {
		t.Fatalf("error should mention invalid min/size; got %q", msg)
	}
}
```
