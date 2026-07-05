## Expected

- HTTP status 501.
- Response body indicates unsupported platform.

## Errors

- No harness error is returned.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 501 {
		t.Fatalf("status = %d, want 501; body=%s", resp.StatusCode, resp.Body)
	}
	if !strings.Contains(resp.Body, "unsupported platform") {
		t.Fatalf("body = %q, want unsupported platform error", resp.Body)
	}
}
```