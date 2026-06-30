## Expected

- One hit is returned with `Kind == "macho"`.
- The hit is streamed to stdout.

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
	if resp.Result == nil || len(resp.Result.Binaries) != 1 {
		t.Fatalf("expected one hit, got %#v", resp.Result)
	}
	hit := resp.Result.Binaries[0]
	if hit.Kind != "macho" {
		t.Fatalf("expected macho kind, got %#v", hit)
	}
	if hit.RepoName != "macho-app" {
		t.Fatalf("expected repo name macho-app, got %#v", hit)
	}
	if !strings.Contains(resp.Stdout, "macho") || !strings.Contains(resp.Stdout, "~/Projects/macho-app/bin/macho-app") {
		t.Fatalf("expected streamed macho line, got:\n%s", resp.Stdout)
	}
}
```
