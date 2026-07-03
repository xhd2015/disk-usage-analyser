## Expected

- `resp.Result.NamedHits` has exactly two entries.
- One has `Name="node_modules"`, size 6; the other has `Name="vendor"`, size 5.
- `resp.Result.TotalSize` is 11.
- Summary: `Found 0 binaries, 2 named entries, total 11 B`.

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
	seen := map[string]int64{}
	for _, hit := range resp.Result.NamedHits {
		seen[hit.Name] = hit.Size
	}
	if seen["node_modules"] != 6 {
		t.Fatalf("expected node_modules size=6, got %d", seen["node_modules"])
	}
	if seen["vendor"] != 5 {
		t.Fatalf("expected vendor size=5, got %d", seen["vendor"])
	}
	if resp.Result.TotalSize != 11 {
		t.Fatalf("expected TotalSize=11, got %d", resp.Result.TotalSize)
	}
	if !strings.Contains(resp.Stdout, "name:node_modules") {
		t.Fatalf("stdout missing name:node_modules:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "name:vendor") {
		t.Fatalf("stdout missing name:vendor:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 2 named entries, total 11 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}
```
