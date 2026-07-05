## Expected

- At least two `named` events and one `scan_complete`.
- Last `named` appears before `scan_complete`.
- `named_enriched` may appear before `scan_complete` (pipelined enrichment).
- `scan_complete` appears before `done`.
- `named_enriched` appears before `done`.
- SSE output contains `event: scan_complete` and `event: named_enriched`.

## Errors

- No harness error is returned.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedJSON) < 2 {
		t.Fatalf("expected at least 2 named hits, got %d", len(resp.NamedJSON))
	}
	if !resp.HasScanComplete {
		t.Fatalf("expected scan_complete event, got events: %v", resp.EventTypes)
	}
	if !resp.LastNamedBeforeScanComplete {
		t.Fatalf("expected last named before scan_complete, got events: %v", resp.EventTypes)
	}
	if !resp.ScanCompleteBeforeDone {
		t.Fatalf("expected scan_complete before done, got events: %v", resp.EventTypes)
	}
	if !resp.NamedEnrichedBeforeDone {
		t.Fatalf("expected named_enriched before done, got events: %v", resp.EventTypes)
	}
	if !strings.Contains(resp.SSEOutput, "event: scan_complete") {
		t.Fatal("expected SSE output to contain 'event: scan_complete'")
	}
	if !strings.Contains(resp.SSEOutput, "event: named_enriched") {
		t.Fatal("expected SSE output to contain 'event: named_enriched'")
	}
}
```