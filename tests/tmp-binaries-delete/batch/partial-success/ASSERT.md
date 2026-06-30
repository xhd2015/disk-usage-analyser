## Expected

- Exactly one path in `deleted` (the binary file).
- Exactly one path in `failed` (the directory).
- `freedSize` reflects only the deleted binary.
- `bin/ok` no longer exists; repo directory remains.

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
	if len(resp.Result.Deleted) != 1 {
		t.Fatalf("deleted = %#v, want 1", resp.Result.Deleted)
	}
	if len(resp.Result.Failed) != 1 {
		t.Fatalf("failed = %#v, want 1", resp.Result.Failed)
	}
	if !strings.Contains(resp.Result.Deleted[0], "bin/ok") {
		t.Fatalf("deleted path = %q", resp.Result.Deleted[0])
	}
	if !strings.Contains(resp.Result.Failed[0].Path, "Projects/partial") {
		t.Fatalf("failed path = %q", resp.Result.Failed[0].Path)
	}
	if resp.Result.FreedSize <= 0 {
		t.Fatalf("freedSize = %d, want > 0", resp.Result.FreedSize)
	}
	if resp.FileExists[resp.Result.Deleted[0]] {
		t.Fatal("deleted binary should not exist")
	}
	if !resp.FileExists["~/Projects/partial"] {
		t.Fatal("repo directory should remain")
	}
}
```