## Expected

- Exactly one `repo` row for `foo` with `size > 0` (main checkout sizing).
- Exactly one `worktree` hit for the linked checkout at `~/Projects/foo-wt`.
- Linked hit has `isMain=false` and `head` is `feature`.
- No `worktree` event has `isMain=true`.
- `repoName` is `foo` for repo row and linked hit.

## Errors

- No harness error is returned.

```go
import (
	"strings"

	"disk-usage-analyser/server"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var repoRow *WorktreeRepoRow
	for i := range resp.RepoRows {
		if resp.RepoRows[i].RepoName == "foo" {
			repoRow = &resp.RepoRows[i]
			break
		}
	}
	if repoRow == nil {
		t.Fatalf("expected repo row for foo, got %#v", resp.RepoRows)
	}
	if repoRow.Size <= 0 {
		t.Fatalf("foo repo row should have main checkout size > 0: %#v", repoRow)
	}

	var linked *server.WorktreeHit
	for i := range resp.Worktrees {
		if resp.Worktrees[i].IsMain {
			t.Fatalf("main checkout must not appear as worktree event: %#v", resp.Worktrees[i])
		}
		if strings.Contains(resp.Worktrees[i].Path, "foo-wt") {
			linked = &resp.Worktrees[i]
		}
	}
	if linked == nil {
		t.Fatalf("expected linked worktree hit for foo-wt, got %#v", resp.Worktrees)
	}
	if linked.IsMain {
		t.Fatalf("linked hit must have isMain=false: %#v", linked)
	}
	if linked.Head != "feature" {
		t.Fatalf("linked head = %q, want feature", linked.Head)
	}
	if linked.RepoName != "foo" {
		t.Fatalf("linked repoName = %q, want foo", linked.RepoName)
	}
	if len(resp.Worktrees) != 1 {
		t.Fatalf("expected exactly 1 linked worktree event, got %d: %#v", len(resp.Worktrees), resp.Worktrees)
	}
}
```