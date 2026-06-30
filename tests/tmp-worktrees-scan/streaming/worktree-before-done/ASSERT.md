## Expected

- Parsed event order includes at least one `worktree` event before the first `done` event.
- SSE output contains `event: worktree` and `event: done`.
- Each `data:` line after an event is valid JSON.

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
	if !resp.WorktreeBeforeDone {
		t.Fatalf("expected worktree event before done, got events: %v", resp.EventTypes)
	}
	if !strings.Contains(resp.SSEOutput, "event: worktree") {
		t.Fatal("expected SSE output to contain 'event: worktree'")
	}
	if !strings.Contains(resp.SSEOutput, "event: done") {
		t.Fatal("expected SSE output to contain 'event: done'")
	}
	for _, line := range strings.Split(resp.SSEOutput, "\n") {
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var v interface{}
			if err := json.Unmarshal([]byte(data), &v); err != nil {
				t.Fatalf("invalid JSON in SSE data: %s", err)
			}
		}
	}
}
```