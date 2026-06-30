## Expected

- Exactly one `repo` row for `solo` with `size > 0` and `fileCount > 0`.
- Zero `worktree` events are parsed from SSE.
- No worktree hit has `isMain=true`.

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
	var solo *WorktreeRepoRow
	for i := range resp.RepoRows {
		if resp.RepoRows[i].RepoName == "solo" {
			solo = &resp.RepoRows[i]
			break
		}
	}
	if solo == nil {
		t.Fatalf("expected repo row for solo, got %#v", resp.RepoRows)
	}
	if solo.Size <= 0 || solo.FileCount <= 0 || solo.SizeHuman == "" {
		t.Fatalf("solo repo row should have main checkout size: %#v", solo)
	}
	if len(resp.Worktrees) != 0 {
		t.Fatalf("expected zero worktree events for main-only repo, got %d: %#v", len(resp.Worktrees), resp.Worktrees)
	}
	for _, hit := range resp.Worktrees {
		if hit.IsMain {
			t.Fatalf("main checkout must not appear as worktree event: %#v", hit)
		}
	}
	if !strings.Contains(resp.SSEOutput, "event: repo") {
		t.Fatal("expected SSE output to contain event: repo")
	}
	if strings.Contains(resp.SSEOutput, "event: worktree") {
		t.Fatal("main-only repo must not emit worktree events")
	}
}
```