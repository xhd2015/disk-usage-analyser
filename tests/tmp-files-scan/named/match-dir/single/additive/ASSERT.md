## Expected

- `resp.Result.NamedHits` has exactly one entry (node_modules, 6 bytes).
- `resp.Result.Binaries` has exactly one entry (macho, 104 bytes).
- `resp.Result.TotalSize` is 110.
- Human stdout contains both a `name:node_modules` line and a `macho` line.
- Summary: `Found 1 binaries, 1 named entries, total 110 B`.

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
	if len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected 1 binary hit, got %d: %#v", len(resp.Result.Binaries), resp.Result.Binaries)
	}
	named := resp.Result.NamedHits[0]
	if named.Name != "node_modules" || named.Size != 6 {
		t.Fatalf("wrong named hit: %#v", named)
	}
	bin := resp.Result.Binaries[0]
	if bin.Kind != "macho" || bin.Size != 104 {
		t.Fatalf("wrong binary hit: %#v", bin)
	}
	if resp.Result.TotalSize != 110 {
		t.Fatalf("expected TotalSize=110, got %d", resp.Result.TotalSize)
	}
	if !strings.Contains(resp.Stdout, "name:node_modules") {
		t.Fatalf("stdout missing name:node_modules:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "macho") {
		t.Fatalf("stdout missing macho:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 1 binaries, 1 named entries, total 110 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
