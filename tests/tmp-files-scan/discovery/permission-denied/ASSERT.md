## Expected

- The readable repo hit is reported.
- The unreadable repo does not abort the command.
- The summary counts only readable hits.

## Side Effects

- Temporary permissions are restored by cleanup.

## Errors

- Permission denial is handled as a partial scan condition, not a fatal error.

## Exit Code

- 0

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
		t.Fatalf("permission-denied scan should be graceful, got: %v", resp.Err)
	}
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one partial hit, got %#v", resp.Result)
	}
	if !strings.Contains(resp.Stdout, "bin/ok") {
		t.Fatalf("readable hit missing:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "bin/locked") {
		t.Fatalf("locked hit should not be reported:\n%s", resp.Stdout)
	}
}
```
