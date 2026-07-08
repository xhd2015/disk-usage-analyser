## Expected

- Xcode `location` SSE event carries non-empty `runtimeItems` when mock runner returns fixture.
- At least one item has Type containing `iOS 18.5`.

## Side Effects

- None (httptest SSE read).

## Errors

- Fails if xcode location found but runtimeItems empty.

## Exit Code

- Test passes when expectations match.

```go
import (
	"strings"
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
		t.Fatalf("expected runtimeItems on Xcode location event, SSE snippet: %s", snippet)
	}
	found := false
	for _, item := range resp.RuntimeItems {
		if strings.Contains(item.Type, "18.5") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected runtime item for iOS 18.5, got %+v", resp.RuntimeItems)
	}
}
```