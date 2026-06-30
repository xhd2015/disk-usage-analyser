## Expected

- At least two `repo` rows (`resp.Repos >= 2`) for `alpha` and `beta`.
- Each repo row has `size > 0` and `repoPath` containing its `repoName`.
- Zero `worktree` events (main-only repos, no linked checkouts).

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
	if resp.Repos < 2 {
		t.Fatalf("expected at least 2 repo events, got %d", resp.Repos)
	}
	seen := map[string]*WorktreeRepoRow{}
	for i := range resp.RepoRows {
		seen[resp.RepoRows[i].RepoName] = &resp.RepoRows[i]
	}
	for _, name := range []string{"alpha", "beta"} {
		row, ok := seen[name]
		if !ok {
			t.Fatalf("missing repo row for %q: %#v", name, resp.RepoRows)
		}
		if row.Size <= 0 {
			t.Fatalf("repo row %q should have main checkout size > 0: %#v", name, row)
		}
		if !strings.Contains(row.RepoPath, name) {
			t.Fatalf("repoPath %q does not match repoName %q", row.RepoPath, name)
		}
	}
	if len(resp.Worktrees) != 0 {
		t.Fatalf("expected zero worktree events for main-only repos, got %d", len(resp.Worktrees))
	}
}
```