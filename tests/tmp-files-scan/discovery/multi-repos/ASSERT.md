## Expected

- Two hits are returned.
- Repo names are `alpha` and `beta`.
- Each hit has the matching repo path.

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
	if resp.Result == nil || len(resp.Result.Binaries) != 2 {
		t.Fatalf("expected two hits, got %#v", resp.Result)
	}
	seen := map[string]bool{}
	for _, hit := range resp.Result.Binaries {
		seen[hit.RepoName] = true
		if !strings.Contains(hit.RepoPath, hit.RepoName) {
			t.Fatalf("repo path %q does not match repo name %q", hit.RepoPath, hit.RepoName)
		}
	}
	if !seen["alpha"] || !seen["beta"] {
		t.Fatalf("missing repo names in hits: %#v", resp.Result.Binaries)
	}
	if !strings.Contains(resp.Stdout, "Found 2 binaries") {
		t.Fatalf("missing two-hit summary:\n%s", resp.Stdout)
	}
}
```
