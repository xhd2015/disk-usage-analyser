## Expected

- First `repo` event index is less than the first `worktree` event index.
- SSE output contains both `event: repo` and `event: worktree`.

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
	if !resp.RepoBeforeWorktree {
		t.Fatalf("expected repo event before worktree, got events: %v", resp.EventTypes)
	}
	if !strings.Contains(resp.SSEOutput, "event: repo") {
		t.Fatal("expected SSE output to contain 'event: repo'")
	}
}
```