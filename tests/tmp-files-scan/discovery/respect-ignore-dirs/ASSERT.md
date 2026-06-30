## Expected

- Only `bin/keep` is reported.
- Ignored directory names do not appear in stdout.

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
		t.Fatalf("expected only kept hit, got %#v", resp.Result)
	}
	if !strings.Contains(resp.Stdout, "bin/keep") {
		t.Fatalf("kept binary missing from stdout:\n%s", resp.Stdout)
	}
	for _, forbidden := range []string{"vendor/tool", "node_modules", ".venv"} {
		if strings.Contains(resp.Stdout, forbidden) {
			t.Fatalf("ignored path %q appeared in stdout:\n%s", forbidden, resp.Stdout)
		}
	}
}
```
