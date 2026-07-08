## Expected

- Non-zero exit code.
- A non-nil error describes the invalid path.

## Exit Code

- Non-zero (implementation may use 2, matching `analyse`)

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected non-zero exit, got 0 (err=%v)", resp.Err)
	}
	if resp.Err == nil {
		t.Fatal("expected path error, got nil")
	}
}
```