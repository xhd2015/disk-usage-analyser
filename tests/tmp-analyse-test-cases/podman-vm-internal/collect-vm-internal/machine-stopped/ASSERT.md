## Expected

- Machine not running → nil `VmInternal`, no error.

## Side Effects

- None (mock runner).

## Errors

- No SSH attempted; graceful nil.

## Exit Code

- Test passes when vmInternal omitted.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatal("expected graceful skip when machine stopped")
	}
	if resp.VmInternal != nil {
		t.Fatalf("expected nil VmInternal when machine stopped, got %+v", resp.VmInternal)
	}
}
```