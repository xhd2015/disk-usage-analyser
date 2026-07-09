## Expected

- Non-zero exit code (unknown / unsupported flag; no backward-compat alias).
- Prefer an error that mentions `threshold`, unknown flag, or similar.

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
		t.Fatalf("expected non-zero exit for --threshold, got 0 (err=%v stdout=%q)", resp.Err, resp.Stdout)
	}
	// Prefer a parse/flag error; if the flag is ignored and 1B is treated as PATH, that is also a failure mode we reject by requiring Err != nil.
	if resp.Err == nil {
		t.Fatal("expected error for removed --threshold flag (no silent accept)")
	}
	msg := strings.ToLower(resp.Err.Error())
	// Accept flag-parse errors or path errors from mis-parse — still must not succeed.
	_ = msg
}
```
