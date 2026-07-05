## Expected

- At least one `named` event for `vendor`.
- No `named_enriched` events.
- `scan_complete` appears before `done`.
- SSE output does not contain `event: named_enriched`.

## Errors

- No harness error is returned.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedJSON) == 0 {
		t.Fatal("expected at least one named vendor hit")
	}
	if resp.HasNamedEnriched || len(resp.NamedEnrichedJSON) > 0 {
		t.Fatalf("vendor scan must not emit named_enriched, got %d events", len(resp.NamedEnrichedJSON))
	}
	if !resp.HasScanComplete {
		t.Fatalf("expected scan_complete for vendor scan, got events: %v", resp.EventTypes)
	}
	if !resp.ScanCompleteBeforeDone {
		t.Fatalf("expected scan_complete before done for vendor scan, got events: %v", resp.EventTypes)
	}
	if strings.Contains(resp.SSEOutput, "event: named_enriched") {
		t.Fatal("vendor SSE output must not contain 'event: named_enriched'")
	}
}
```