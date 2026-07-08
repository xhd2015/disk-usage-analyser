## Expected

- Command failure: empty runtimeItems, nil error (graceful).

## Side Effects

- None.

## Errors

- Runner error is swallowed; no error returned to caller.

## Exit Code

- Test passes when expectations match.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatal("expected graceful handling when command fails")
	}
	if len(resp.RuntimeItems) != 0 {
		t.Fatalf("expected empty slice, got %d", len(resp.RuntimeItems))
	}
}
```