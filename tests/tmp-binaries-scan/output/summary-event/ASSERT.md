## Expected

- SSE contains `event: summary` before `event: done`.
- Summary reports `Binaries=2` and `Repos=2`.
- `TotalSize` equals sum of hit sizes; `TotalHuman` is non-empty.

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
	if !strings.Contains(resp.SSEOutput, "event: summary") {
		t.Fatal("expected summary event")
	}
	if resp.Summary == nil {
		t.Fatal("expected parsed summary payload")
	}
	if resp.Summary.Binaries != 2 || resp.Summary.Repos != 2 {
		t.Fatalf("summary counts = %#v, want 2 binaries and 2 repos", resp.Summary)
	}
	var total int64
	for _, hit := range resp.Binaries {
		total += hit.Size
	}
	if resp.Summary.TotalSize != total {
		t.Fatalf("summary total %d != hit sum %d", resp.Summary.TotalSize, total)
	}
	if resp.Summary.TotalHuman == "" {
		t.Fatal("expected non-empty TotalHuman")
	}
}
```
