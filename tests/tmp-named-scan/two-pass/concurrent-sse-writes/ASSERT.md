## Expected

- Exactly three `named` events and three `named_enriched` events (one enrichment per hit).
- Summary `namedHits` count matches hit count.
- No `server_error` events in the SSE stream.
- `scan_complete` and `done` events are present.

## Errors

- No harness error is returned.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const wantHits = 3
	if len(resp.NamedJSON) != wantHits {
		t.Fatalf("named event count = %d, want %d", len(resp.NamedJSON), wantHits)
	}
	if len(resp.NamedEnrichedJSON) != wantHits {
		t.Fatalf("named_enriched event count = %d, want %d", len(resp.NamedEnrichedJSON), wantHits)
	}
	if resp.Summary == nil {
		t.Fatal("expected summary event")
	}
	if resp.Summary.NamedHits != wantHits {
		t.Fatalf("summary namedHits = %d, want %d", resp.Summary.NamedHits, wantHits)
	}
	if resp.HasServerError {
		t.Fatalf("expected no server_error events, got events: %v", resp.EventTypes)
	}
	if strings.Contains(resp.SSEOutput, "event: server_error") {
		t.Fatal("SSE output must not contain 'event: server_error'")
	}
	if !resp.HasScanComplete {
		t.Fatalf("expected scan_complete event, got events: %v", resp.EventTypes)
	}
	if !resp.ScanCompleteBeforeDone {
		t.Fatalf("expected scan_complete before done, got events: %v", resp.EventTypes)
	}
}
```