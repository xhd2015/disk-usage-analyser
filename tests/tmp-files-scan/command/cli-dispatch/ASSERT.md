## Expected

- `run.RunWithOptions` returns nil.
- The injected server starter is not called.
- Stdout contains the binary hit and summary.

## Side Effects

- No HTTP server is started.

## Errors

- No error is returned.

## Exit Code

- Successful command return.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected dispatch error: %v", resp.Err)
	}
	if resp.ServerWasStarted {
		t.Fatalf("tmp-files dispatch started web server")
	}
	if !strings.Contains(resp.Stdout, "~/Projects/app/bin/app") {
		t.Fatalf("expected app binary in stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 1 binaries, total ") {
		t.Fatalf("expected summary in stdout:\n%s", resp.Stdout)
	}
}
```
