## Expected

- Zero worktree hits are parsed from SSE.
- Summary reports `Worktrees=0` and `Repos=0` when present.
- `done` event is still emitted.

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
	if len(resp.Worktrees) != 0 {
		t.Fatalf("expected zero worktree hits, got %d: %#v", len(resp.Worktrees), resp.Worktrees)
	}
	if resp.Summary != nil && (resp.Summary.Worktrees != 0 || resp.Summary.Repos != 0) {
		t.Fatalf("expected zero summary counts, got %#v", resp.Summary)
	}
	if !strings.Contains(resp.SSEOutput, "event: done") {
		t.Fatal("expected done event even with no repos")
	}
}
```