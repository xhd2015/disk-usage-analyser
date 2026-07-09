## Expected

- Non-zero exit code.
- A non-nil error describes the missing path / file not found.

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
		t.Fatalf("expected non-zero exit for missing path, got 0 (err=%v)", resp.Err)
	}
	if resp.Err == nil && resp.Stderr == "" {
		t.Fatal("expected path error or stderr message, got neither")
	}
	msg := ""
	if resp.Err != nil {
		msg = resp.Err.Error()
	}
	msg += " " + resp.Stderr
	lower := strings.ToLower(msg)
	// Accept common phrasings: no such file, not found, does not exist, stat failure.
	ok := strings.Contains(lower, "no such") ||
		strings.Contains(lower, "not found") ||
		strings.Contains(lower, "does not exist") ||
		strings.Contains(lower, "exist") ||
		strings.Contains(lower, "stat") ||
		(req.TargetPath != "" && strings.Contains(msg, req.TargetPath))
	if !ok {
		t.Fatalf("error should indicate missing path, got: err=%v stderr=%q", resp.Err, resp.Stderr)
	}
}
```
