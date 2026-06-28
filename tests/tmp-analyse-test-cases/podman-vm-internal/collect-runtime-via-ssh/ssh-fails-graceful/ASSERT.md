## Expected

- Inner SSH failure → empty `runtimeItems`, no collection error.

## Side Effects

- None (mock runner).

## Errors

- Graceful empty result; scan would continue.

## Exit Code

- Test passes when runtime items empty without error.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatal("expected graceful handling when inner system df fails")
	}
	if len(resp.RuntimeItems) != 0 {
		t.Fatalf("expected empty runtime items, got %d", len(resp.RuntimeItems))
	}
}
```