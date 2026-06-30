## Expected

- Scanned path is listed in `failed`.
- Error mentions not found or missing file.
- `deleted` is empty.
- `freedSize` is zero.

## Errors

- No harness error is returned.

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Result.Deleted) != 0 {
		t.Fatalf("deleted = %#v, want empty", resp.Result.Deleted)
	}
	if len(resp.Result.Failed) != 1 {
		t.Fatalf("failed = %#v, want 1", resp.Result.Failed)
	}
	errMsg := strings.ToLower(resp.Result.Failed[0].Error)
	if !strings.Contains(errMsg, "not found") &&
		!strings.Contains(errMsg, "no such") &&
		!strings.Contains(errMsg, "exist") {
		t.Fatalf("expected not-found error, got %q", resp.Result.Failed[0].Error)
	}
	if resp.Result.FreedSize != 0 {
		t.Fatalf("freedSize = %d, want 0", resp.Result.FreedSize)
	}
}
```