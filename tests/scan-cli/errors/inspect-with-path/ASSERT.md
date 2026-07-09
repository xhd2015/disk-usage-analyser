## Expected

- Non-zero exit code.
- Clear error that inspect mode does not accept a positional PATH (or unexpected argument).

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
		t.Fatalf("expected non-zero exit when combining --inspect with PATH, got 0")
	}
	if resp.Err == nil {
		t.Fatal("expected error for --inspect + PATH, got nil")
	}
	msg := strings.ToLower(resp.Err.Error())
	// Accept messages about unexpected arg, mutual exclusion, or inspect+path conflict.
	if !strings.Contains(msg, "inspect") &&
		!strings.Contains(msg, "unexpected") &&
		!strings.Contains(msg, "path") &&
		!strings.Contains(msg, "argument") {
		t.Fatalf("error should explain inspect/path conflict; got %q", msg)
	}
}
```
