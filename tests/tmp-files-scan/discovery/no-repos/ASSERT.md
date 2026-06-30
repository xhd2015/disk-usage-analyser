## Expected

- `Repos` is 0.
- No binary hits are returned.
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
	if resp.Result == nil {
		t.Fatalf("expected result")
	}
	if resp.Result.Repos != 0 {
		t.Fatalf("expected zero repos, got %d", resp.Result.Repos)
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected zero binaries, got %#v", resp.Result.Binaries)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries") {
		t.Fatalf("missing zero summary:\n%s", resp.Stdout)
	}
}
```
