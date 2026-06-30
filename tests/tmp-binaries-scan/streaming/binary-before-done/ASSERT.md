## Expected

- Parsed event order includes at least one `binary` event before `done`.
- SSE output contains `event: binary` and `event: done`.
- Each `data:` line is valid JSON.

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
	if !resp.BinaryBeforeDone {
		t.Fatalf("expected binary event before done, got events: %v", resp.EventTypes)
	}
	if !strings.Contains(resp.SSEOutput, "event: binary") {
		t.Fatal("expected SSE output to contain 'event: binary'")
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
