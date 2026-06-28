## Expected

- Non-darwin platform → nil `VmInternal`, no error, no SSH calls required.

## Side Effects

- None.

## Errors

- Graceful platform skip.

## Exit Code

- Test passes when vmInternal omitted on non-darwin.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.CollectFailed {
		t.Fatal("expected graceful platform skip")
	}
	if resp.VmInternal != nil {
		t.Fatalf("expected nil VmInternal on non-darwin, got %+v", resp.VmInternal)
	}
}
```