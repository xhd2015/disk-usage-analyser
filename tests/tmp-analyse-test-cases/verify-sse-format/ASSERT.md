## Expected
- SSE output contains at least one "event: location" line
- SSE output contains "event: summary" line
- SSE output contains "event: done" line
- Each data line after an event contains valid JSON

## Errors
- If the SSE output is empty, fail with message "expected SSE events, got empty output"

```go
import (
	"encoding/json"
	"strings"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SSEOutput == "" {
		t.Fatal("expected SSE events, got empty output")
	}
	if !strings.Contains(resp.SSEOutput, "event: location") {
		t.Fatal("expected SSE output to contain 'event: location'")
	}
	if !strings.Contains(resp.SSEOutput, "event: summary") {
		t.Fatal("expected SSE output to contain 'event: summary'")
	}
	if !strings.Contains(resp.SSEOutput, "event: done") {
		t.Fatal("expected SSE output to contain 'event: done'")
	}

	lines := strings.Split(resp.SSEOutput, "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var v interface{}
			if err := json.Unmarshal([]byte(data), &v); err != nil {
				t.Fatalf("line %d: invalid JSON in SSE data: %s", i+1, err)
			}
		}
	}
}
```
