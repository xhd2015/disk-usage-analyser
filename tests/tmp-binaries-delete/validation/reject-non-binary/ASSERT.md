## Expected

- Scanned path appears in `failed` with a non-empty error.
- `deleted` is empty.
- File still exists on disk (overwrite content remains).

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
	if !strings.Contains(resp.Result.Failed[0].Path, "bin/good") {
		t.Fatalf("failed path = %q", resp.Result.Failed[0].Path)
	}
	if resp.Result.Failed[0].Error == "" {
		t.Fatal("expected failure error message")
	}
}
```