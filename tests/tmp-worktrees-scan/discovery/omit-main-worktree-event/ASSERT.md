## Expected

- Raw SSE contains at least one `event: worktree` payload.
- No worktree event JSON contains `"isMain":true`.
- Repo row for `omit-main` has `size > 0` (main checkout on repo event).

## Errors

- No harness error is returned.

```go
import (
	"encoding/json"
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(resp.SSEOutput, "event: worktree") {
		t.Fatal("expected at least one worktree event for linked checkout")
	}

	var currentEvent string
	for _, line := range strings.Split(resp.SSEOutput, "\n") {
		if strings.HasPrefix(line, "event: ") {
			currentEvent = strings.TrimPrefix(line, "event: ")
		}
		if strings.HasPrefix(line, "data: ") && currentEvent == "worktree" {
			data := strings.TrimPrefix(line, "data: ")
			var payload map[string]interface{}
			if err := json.Unmarshal([]byte(data), &payload); err != nil {
				t.Fatalf("invalid worktree JSON: %v", err)
			}
			if isMain, ok := payload["isMain"].(bool); ok && isMain {
				t.Fatalf("worktree event must not have isMain=true: %s", data)
			}
		}
	}

	for i := range resp.Worktrees {
		if resp.Worktrees[i].IsMain {
			t.Fatalf("parsed worktree hit must not have isMain=true: %#v", resp.Worktrees[i])
		}
	}

	var repoRow *WorktreeRepoRow
	for i := range resp.RepoRows {
		if resp.RepoRows[i].RepoName == "omit-main" {
			repoRow = &resp.RepoRows[i]
			break
		}
	}
	if repoRow == nil || repoRow.Size <= 0 {
		t.Fatalf("expected repo row with main checkout size for omit-main: %#v", repoRow)
	}
}
```