## Expected

- `resp.Result.NamedHits` has exactly one entry (only the selected subtree).
- NamedHit `Size` is 5 (from the 5-byte file in selected).
- `resp.Result.Binaries` is empty.
- stdout does not mention the unselected repo.

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
	if len(resp.Result.NamedHits) != 1 {
		t.Fatalf("expected 1 named hit, got %d: %#v", len(resp.Result.NamedHits), resp.Result.NamedHits)
	}
	hit := resp.Result.NamedHits[0]
	if hit.Size != 5 {
		t.Fatalf("expected Size=5, got %d", hit.Size)
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected no binary hits, got %#v", resp.Result.Binaries)
	}
	if strings.Contains(resp.Stdout, "unselected") {
		t.Fatalf("unselected path leaked into stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 1 named entries, total 5 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
