## Expected

- `resp.Result.NamedHits` has exactly one entry (`node_modules`, 6 bytes).
- `resp.Result.Binaries` is empty (`vendor/tool` is skipped as ignored).
- `resp.Result.TotalSize` is 6.
- stdout does not mention `vendor` or `macho`.
- Summary: `Found 0 binaries, 1 named entries, total 6 B`.

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
	if hit.Name != "node_modules" || hit.Size != 6 {
		t.Fatalf("wrong named hit: %#v", hit)
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected 0 binary hits, got %d: %#v", len(resp.Result.Binaries), resp.Result.Binaries)
	}
	if resp.Result.TotalSize != 6 {
		t.Fatalf("expected TotalSize=6, got %d", resp.Result.TotalSize)
	}
	if strings.Contains(resp.Stdout, "vendor") {
		t.Fatalf("vendor appeared in stdout:\n%s", resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "macho") {
		t.Fatalf("macho appeared in stdout:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 1 named entries, total 6 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}
```
