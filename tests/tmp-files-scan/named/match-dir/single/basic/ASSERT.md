## Expected

- `resp.Result.NamedHits` has exactly one entry.
- NamedHit fields: `Path` ends with `node_modules`, `Name` is `"node_modules"`, `Size` is 6, `SizeHuman` is `"6 B"`, `RepoName` is `"app"`.
- Human stdout contains `name:node_modules` and tilde-expanded paths.
- Summary line: `Found 0 binaries, 1 named entries, total 6 B`.
- `resp.Result.Binaries` is empty.
- `resp.Result.TotalSize` is 6.

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
	if hit.Size != 6 {
		t.Fatalf("expected Size=6, got %d", hit.Size)
	}
	if hit.SizeHuman != "6 B" {
		t.Fatalf("expected SizeHuman='6 B', got %q", hit.SizeHuman)
	}
	if hit.RepoName != "app" {
		t.Fatalf("expected RepoName=app, got %q", hit.RepoName)
	}
	if !strings.Contains(hit.Path, "~/Projects/app/node_modules") {
		t.Fatalf("expected path containing ~/Projects/app/node_modules, got %q", hit.Path)
	}
	if !strings.Contains(hit.RepoPath, "~/Projects/app") {
		t.Fatalf("expected repoPath containing ~/Projects/app, got %q", hit.RepoPath)
	}
	if len(resp.Result.Binaries) != 0 {
		t.Fatalf("expected no binary hits, got %d: %#v", len(resp.Result.Binaries), resp.Result.Binaries)
	}
	if resp.Result.TotalSize != 6 {
		t.Fatalf("expected TotalSize=6, got %d", resp.Result.TotalSize)
	}
	if !strings.Contains(resp.Stdout, "name:node_modules") {
		t.Fatalf("stdout missing name:node_modules:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "6 B") {
		t.Fatalf("stdout missing '6 B':\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "Found 0 binaries, 1 named entries, total 6 B") {
		t.Fatalf("bad summary line:\n%s", resp.Stdout)
	}
}

```
