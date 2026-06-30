## Expected

- Two binary hits are returned.
- Repo names are `alpha` and `beta`.
- Each hit's `repoPath` contains its `repoName`.

## Errors

- No harness error is returned.

```go
import (
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Binaries) != 2 {
		t.Fatalf("expected two hits, got %#v", resp.Binaries)
	}
	seen := map[string]bool{}
	for _, hit := range resp.Binaries {
		seen[hit.RepoName] = true
		if !strings.Contains(hit.RepoPath, hit.RepoName) {
			t.Fatalf("repoPath %q does not match repoName %q", hit.RepoPath, hit.RepoName)
		}
	}
	if !seen["alpha"] || !seen["beta"] {
		t.Fatalf("missing repo names: %#v", resp.Binaries)
	}
}
```
