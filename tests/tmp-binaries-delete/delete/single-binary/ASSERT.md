## Expected

- HTTP status 200.
- Exactly one path in `deleted`.
- `failed` is empty.
- `freedSize` equals the binary file size from scan.
- File no longer exists on disk.

## Errors

- No harness error is returned.

```go
import (
	"os"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200; body=%s", resp.StatusCode, resp.Body)
	}
	if resp.Result == nil {
		t.Fatal("expected delete result")
	}
	if len(resp.Result.Deleted) != 1 {
		t.Fatalf("deleted = %#v, want 1 path", resp.Result.Deleted)
	}
	if len(resp.Result.Failed) != 0 {
		t.Fatalf("failed = %#v, want empty", resp.Result.Failed)
	}
	if resp.Result.FreedSize <= 0 {
		t.Fatalf("freedSize = %d, want > 0", resp.Result.FreedSize)
	}
	for _, p := range resp.Result.Deleted {
		if resp.FileExists[p] {
			t.Fatalf("file still exists after delete: %s", p)
		}
	}
}
```
