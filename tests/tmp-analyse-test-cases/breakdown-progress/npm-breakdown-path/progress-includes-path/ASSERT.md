## Expected
- npm progress event includes breakdownPath for dynamic subdir row.

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
	p := resp.ProgressPayload
	if p["breakdownPath"] != "~/.npm/_cacache" { t.Fatalf("expected breakdownPath, got %v", p["breakdownPath"]) }
	if p["breakdownSize"].(int64) != 1024 { t.Fatalf("expected breakdownSize=1024, got %d", p["breakdownSize"]) }
}
```
