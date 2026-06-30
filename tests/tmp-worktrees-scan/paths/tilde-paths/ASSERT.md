## Expected

- At least one `repo` row has `repoPath` starting with `~/`.
- Linked worktree hits (if any) have `path` and `repoPath` starting with `~/`.
- No row exposes the absolute `HomeDir` path from the harness.

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
	if len(resp.RepoRows) == 0 {
		t.Fatal("expected at least one repo row")
	}
	for _, row := range resp.RepoRows {
		if !strings.HasPrefix(row.RepoPath, "~/") {
			t.Fatalf("repoPath %q should start with ~/", row.RepoPath)
		}
		if strings.Contains(row.RepoPath, req.HomeDir) {
			t.Fatalf("repo row exposes absolute home path: %#v", row)
		}
	}
	for _, hit := range resp.Worktrees {
		if !strings.HasPrefix(hit.Path, "~/") {
			t.Fatalf("path %q should start with ~/", hit.Path)
		}
		if !strings.HasPrefix(hit.RepoPath, "~/") {
			t.Fatalf("repoPath %q should start with ~/", hit.RepoPath)
		}
		if strings.Contains(hit.Path, req.HomeDir) || strings.Contains(hit.RepoPath, req.HomeDir) {
			t.Fatalf("hit exposes absolute home path: %#v", hit)
		}
	}
}
```