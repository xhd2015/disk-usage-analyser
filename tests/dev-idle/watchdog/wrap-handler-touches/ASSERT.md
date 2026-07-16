## Expected

- Wrapped handler returns HTTP 200.
- `OnIdle` is not called immediately after the HTTP request (wrap touch resets idle clock).
- After advancing 3s past the request with `Timeout=2s`, `OnIdle` is called exactly once.

## Errors

- No harness error is returned.

```go
import "net/http"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.HTTPStatus != http.StatusOK {
		t.Fatalf("HTTPStatus = %d, want 200", resp.HTTPStatus)
	}
	if resp.OnIdleCountAfterHTTP != 0 {
		t.Fatalf("OnIdleCountAfterHTTP = %d, want 0 immediately after request", resp.OnIdleCountAfterHTTP)
	}
	if resp.OnIdleCount != 1 {
		t.Fatalf("OnIdleCount = %d, want 1 after idle period following wrapped request", resp.OnIdleCount)
	}
}
```