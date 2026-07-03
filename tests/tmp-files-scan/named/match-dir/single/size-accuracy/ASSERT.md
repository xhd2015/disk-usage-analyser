## Expected

- `resp.Result.NamedHits` has exactly one entry.
- NamedHit `Size` is exactly 600 (100+200+300).
- `SizeHuman` reflects 600 B.
- `resp.Result.TotalSize` is 600.

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
	if hit.Name != "node_modules" {
		t.Fatalf("expected Name=node_modules, got %q", hit.Name)
	}
	if hit.Size != 600 {
		t.Fatalf("expected Size=600, got %d", hit.Size)
	}
	if hit.SizeHuman != "600 B" {
		t.Fatalf("expected SizeHuman='600 B', got %q", hit.SizeHuman)
	}
	if resp.Result.TotalSize != 600 {
		t.Fatalf("expected TotalSize=600, got %d", resp.Result.TotalSize)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 1 named entries, total 600 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
