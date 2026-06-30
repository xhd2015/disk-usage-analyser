## Expected

- Directory path is listed in `failed`.
- Directory still exists.
- Binary `bin/app` still exists.

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
	errMsg := strings.ToLower(resp.Result.Failed[0].Error)
	if errMsg == "" {
		t.Fatal("expected non-empty failure error")
	}
	if !strings.Contains(errMsg, "director") && !strings.Contains(errMsg, "not a file") {
		t.Fatalf("expected directory rejection error, got %q", resp.Result.Failed[0].Error)
	}
	if !resp.FileExists["~/Projects/bad-dir"] {
		t.Fatal("directory should remain")
	}
}
```
