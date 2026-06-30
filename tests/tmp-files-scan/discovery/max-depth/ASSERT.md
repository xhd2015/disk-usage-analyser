## Expected

- Only the shallow repository hit is reported.
- The deep repository hit is absent.

## Side Effects

- None outside the temporary fixture tree.

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
		t.Fatalf("expected only shallow hit, got %#v", resp.Result)
	}
	if !strings.Contains(resp.Stdout, "bin/shallow") {
		t.Fatalf("shallow hit missing:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "bin/deep") {
		t.Fatalf("deep hit should be excluded by max-depth:\n%s", resp.Stdout)
	}
}
```
