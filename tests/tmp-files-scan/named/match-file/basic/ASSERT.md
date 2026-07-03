## Expected

- `resp.Result.NamedHits` has exactly one entry.
- NamedHit fields: `Name="node_modules"`, `Size=5`, `SizeHuman="5 B"`, `Path` ends with `/node_modules`.
- `resp.Result.Binaries` is empty (the file is not a binary).
- Human stdout contains `name:node_modules`.
- Summary: `Found 0 binaries, 1 named entries, total 5 B`.

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
	if hit.SizeHuman != "5 B" {
		t.Fatalf("expected SizeHuman='5 B', got %q", hit.SizeHuman)
	}
	if !strings.Contains(hit.Path, "/node_modules") {
		t.Fatalf("expected path containing /node_modules, got %q", hit.Path)
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
