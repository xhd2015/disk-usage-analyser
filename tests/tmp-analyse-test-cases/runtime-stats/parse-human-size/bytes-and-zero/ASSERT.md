## Expected
- "500 B" -> 500; "0 B" -> 0.

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
	if resp.ParsedBytes != 500 { t.Fatalf("expected 500 bytes, got %d", resp.ParsedBytes) }
	zero, err := server.ParseHumanSize("0 B")
	if err != nil || zero != 0 { t.Fatalf("expected 0 B -> 0, got %d err=%v", zero, err) }
}
```
