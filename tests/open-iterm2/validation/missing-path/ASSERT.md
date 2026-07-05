## Expected

- HTTP status 400.
- Response body mentions path is required.

## Errors

- No harness error is returned.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 400 {
		t.Fatalf("status = %d, want 400; body=%s", resp.StatusCode, resp.Body)
	}
	if !strings.Contains(resp.Body, "path required") {
		t.Fatalf("body = %q, want path required error", resp.Body)
	}
}
```