## Expected

- No hits are returned.
- The outside binary path is absent from stdout.
- Summary says `Found 0 binaries`.

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
	if resp.Result == nil || len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected no outside-repo hits, got %#v", resp.Result)
	}
	if strings.Contains(resp.Stdout, "Downloads/tool") {
		t.Fatalf("outside repo binary appeared in stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries") {
		t.Fatalf("missing zero summary:\n%s", resp.Stdout)
	}
}
```
