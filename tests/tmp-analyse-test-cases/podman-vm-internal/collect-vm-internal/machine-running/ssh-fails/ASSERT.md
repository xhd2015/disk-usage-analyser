## Expected

- SSH error → nil `VmInternal`, no collection error.

## Side Effects

- None (mock runner).

## Errors

- `CollectFailed` must be false; graceful nil result.

## Exit Code

- Test passes when vmInternal omitted without error.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatal("expected graceful handling when SSH fails")
	}
	if resp.VmInternal != nil {
		t.Fatalf("expected nil VmInternal on SSH failure, got %+v", resp.VmInternal)
	}
}
```