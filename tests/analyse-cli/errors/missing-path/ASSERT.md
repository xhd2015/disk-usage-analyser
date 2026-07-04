## Expected

- Exit code 2.
- A non-nil error describes the invalid path.

## Exit Code

- 2

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 2 {
		t.Fatalf("expected exit code 2, got %d (err=%v)", resp.ExitCode, resp.Err)
	}
	if resp.Err == nil {
		t.Fatal("expected usage/path error, got nil")
	}
}
```