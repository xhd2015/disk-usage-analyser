## Expected

- `resp.Result.NamedHits` has exactly two entries.
- Total size is 13 (6 + 7).
- The two hits have distinct paths; neither is a subset of the other (nested hits are separate).
- Summary line: `Found 0 binaries, 2 named entries, total 13 B`.

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
	if len(resp.Result.NamedHits) != 2 {
		t.Fatalf("expected 2 named hits, got %d: %#v", len(resp.Result.NamedHits), resp.Result.NamedHits)
	}
	for _, hit := range resp.Result.NamedHits {
		if hit.Name != "node_modules" {
			t.Fatalf("expected Name=node_modules, got %q", hit.Name)
		}
	}
	// Outer hit (6 bytes, outer.txt only)
	var foundOuter, foundInner bool
	for _, hit := range resp.Result.NamedHits {
		if strings.Contains(hit.Path, "pkg/node_modules") {
			foundInner = true
			if hit.Size != 7 {
				t.Fatalf("expected inner size 7, got %d for %s", hit.Size, hit.Path)
			}
		} else {
			foundOuter = true
			if hit.Size != 6 {
				t.Fatalf("expected outer size 6, got %d for %s", hit.Size, hit.Path)
			}
		}
	}
	if !foundOuter || !foundInner {
		t.Fatalf("missing outer or inner hit: %#v", resp.Result.NamedHits)
	}
	if resp.Result.TotalSize != 13 {
		t.Fatalf("expected TotalSize=13, got %d", resp.Result.TotalSize)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 2 named entries, total 13 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
