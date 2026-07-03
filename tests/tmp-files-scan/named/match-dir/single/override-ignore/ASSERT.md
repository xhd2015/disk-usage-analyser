## Expected

- `resp.Result.NamedHits` has exactly one entry for `node_modules` with Name `"node_modules"` and size 5.
- `resp.Result.Binaries` is empty (no binary detection inside the named dir).
- Summary reports 1 named entry, 0 binaries.

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
	if hit.Size != 5 {
		t.Fatalf("expected Size=5, got %d", hit.Size)
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected no binary hits, got %#v", resp.Result.Binaries)
	}
	if !strings.Contains(resp.Stdout, "name:node_modules") {
		t.Fatalf("stdout missing name:node_modules:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 1 named entries, total 5 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
