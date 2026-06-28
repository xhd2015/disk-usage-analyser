## Expected
- Docker location event carries runtimeItems when mock runner returns fixture.

## Side Effects
- None (pure function or mock CLI).

## Errors
- See leaf scenario for expected error vs graceful-empty behavior.

## Exit Code
- Test passes when expectations match.

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.RuntimeItems) == 0 {
		snippet := resp.SSEOutput
		if len(snippet) > 200 {
			snippet = snippet[:200]
		}
		t.Fatalf("expected runtimeItems on Docker location event, SSE snippet: %s", snippet)
	}
}
```
