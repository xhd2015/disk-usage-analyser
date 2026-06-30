## Expected

- `repo` row for `size-app` has `size > 0` and `fileCount > 0`.
- Linked worktree hit for path containing `size-wt` has `size > 0`, `fileCount > 0`, and non-empty `sizeHuman`.
- No worktree hit has `isMain=true`.

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
		if strings.Contains(resp.RepoRows[i].RepoName, "size-app") {
			repoRow = &resp.RepoRows[i]
			break
		}
	}
	if repoRow == nil || repoRow.Size <= 0 || repoRow.FileCount <= 0 {
		t.Fatalf("main repo row should have size > 0: %#v", repoRow)
	}

	var linked *server.WorktreeHit
	for i := range resp.Worktrees {
		if resp.Worktrees[i].IsMain {
			t.Fatalf("main checkout must not appear as worktree event: %#v", resp.Worktrees[i])
		}
		if strings.Contains(resp.Worktrees[i].Path, "size-wt") {
			linked = &resp.Worktrees[i]
		}
	}
	if linked == nil {
		t.Fatalf("missing linked worktree hit: %#v", resp.Worktrees)
	}
	if linked.Size <= 0 || linked.FileCount <= 0 || linked.SizeHuman == "" {
		t.Fatalf("linked worktree sizing empty: %#v", linked)
	}
}
```