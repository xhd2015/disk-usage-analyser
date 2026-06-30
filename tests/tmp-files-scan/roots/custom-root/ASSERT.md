## Expected

- Only the selected subtree hit is reported.
- The unselected repo does not appear in stdout or result hits.

## Side Effects

- No files outside the explicit root are scanned.

## Errors

- No error is returned.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected scan error: %v", resp.Err)
	}
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one selected binary, got %#v", resp.Result)
	}
	if !strings.Contains(resp.Stdout, "selected-app") {
		t.Fatalf("expected selected hit in stdout:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "unselected-app") {
		t.Fatalf("unselected root leaked into stdout:\n%s", resp.Stdout)
	}
}
```
