## Expected

- Stdout receives more than one write.
- The first write contains a binary hit path.
- The first write does not contain the final `Found` summary.

## Side Effects

- None outside the temporary fixture tree.

## Errors

- No error is returned.

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
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	if resp.WriteCount < 2 {
		t.Fatalf("expected separate writes for streaming hits and summary, got %d writes:\n%s", resp.WriteCount, resp.Stdout)
	}
	if !strings.Contains(resp.FirstWrite, "~/Projects/stream-app/bin/") {
		t.Fatalf("first write should contain a hit line, got %q", resp.FirstWrite)
	}
	if strings.Contains(resp.FirstWrite, "Found ") {
		t.Fatalf("first write includes final summary, output was buffered: %q", resp.FirstWrite)
	}
	if !strings.Contains(resp.Stdout, "Found 2 binaries") {
		t.Fatalf("missing two-hit summary:\n%s", resp.Stdout)
	}
}
```
