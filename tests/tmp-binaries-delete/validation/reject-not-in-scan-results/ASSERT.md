## Expected

- `~/Projects/scope/bin/secret` appears in `failed`.
- Error indicates path was not in scan results (or not allowed).
- `bin/secret` file still exists on disk.
- `bin/scanned` remains untouched.

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
	if len(resp.Result.Failed) != 1 {
		t.Fatalf("failed = %#v, want 1", resp.Result.Failed)
	}
	if !strings.Contains(resp.Result.Failed[0].Path, "bin/secret") {
		t.Fatalf("failed path = %q", resp.Result.Failed[0].Path)
	}
	if resp.Result.Failed[0].Error == "" {
		t.Fatal("expected failure error")
	}
	if !resp.FileExists["~/Projects/scope/bin/secret"] {
		t.Fatal("secret binary should remain on disk")
	}
}
```