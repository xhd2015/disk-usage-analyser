## Expected

- Non-zero exit code when neither PATH nor `--kind` nor `--all-kinds` is provided.
- A non-nil error (or stderr message) indicates PATH is required / missing (or usage),
  and/or that a kind / all-kinds flag is needed.
- Contrast: `explain --kind xcode` without PATH is allowed (see kinds/xcode/no-path);
  `explain --all-kinds` without PATH is allowed (see output/all-kinds-index).

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
		t.Fatalf("expected non-zero exit when PATH is missing and --kind/--all-kinds are unset, got 0 (err=%v stdout=%q)", resp.Err, resp.Stdout)
	}
	// Prefer structured error; also accept message on stderr.
	msg := ""
	if resp.Err != nil {
		msg = resp.Err.Error()
	}
	msg += " " + resp.Stderr + " " + resp.Stdout
	lower := strings.ToLower(msg)
	if !strings.Contains(lower, "path") && !strings.Contains(lower, "required") &&
		!strings.Contains(lower, "usage") && !strings.Contains(lower, "kind") &&
		!strings.Contains(lower, "all-kinds") {
		t.Fatalf("error should mention missing PATH / required flags (or usage), got exit=%d err=%v stderr=%q stdout=%q",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
}
```