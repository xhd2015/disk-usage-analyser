## Expected

- `resp.Result.NamedHits` is empty.
- `resp.Result.Binaries` is empty.
- `resp.Result.TotalSize` is 0.
- Summary: `Found 0 binaries, 0 named entries, total 0 B`.
- No error.

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
		t.Fatal("expected result")
	}
	if len(resp.Result.NamedHits) != 0 {
		t.Fatalf("expected 0 named hits, got %d: %#v", len(resp.Result.NamedHits), resp.Result.NamedHits)
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected 0 binary hits, got %d: %#v", len(resp.Result.Binaries), resp.Result.Binaries)
	}
	if resp.Result.TotalSize != 0 {
		t.Fatalf("expected TotalSize=0, got %d", resp.Result.TotalSize)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 0 named entries, total 0 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
