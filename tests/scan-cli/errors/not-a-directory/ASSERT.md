## Expected

- Non-zero exit code.
- Error message mentions the path is not a directory.

## Exit Code

- Non-zero

```go
import (
	"strings"
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
		t.Fatal("expected not-a-directory error, got nil")
	}
	msg := strings.ToLower(resp.Err.Error())
	if !strings.Contains(msg, "not a directory") && !strings.Contains(msg, "not a dir") {
		t.Fatalf("error should mention not a directory, got: %v", resp.Err)
	}
}
```