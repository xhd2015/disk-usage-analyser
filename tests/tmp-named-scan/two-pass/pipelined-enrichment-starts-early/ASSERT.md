## Expected

- At least two `named` hits and one `named_enriched` event.
- First `named_enriched` event index is strictly less than `scan_complete` event index (enrichment starts during pass 1).

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedJSON) < 2 {
		t.Fatalf("expected at least 2 named hits, got %d", len(resp.NamedJSON))
	}
	if !resp.HasNamedEnriched {
		t.Fatalf("expected named_enriched events, got events: %v", resp.EventTypes)
	}
	if resp.FirstNamedEnrichedIndex < 0 {
		t.Fatalf("expected first named_enriched index, got events: %v", resp.EventTypes)
	}
	if resp.FirstScanCompleteIndex < 0 {
		t.Fatalf("expected scan_complete index, got events: %v", resp.EventTypes)
	}
	if resp.FirstNamedEnrichedIndex >= resp.FirstScanCompleteIndex {
		t.Fatalf("expected first named_enriched (idx %d) before scan_complete (idx %d), got events: %v",
			resp.FirstNamedEnrichedIndex, resp.FirstScanCompleteIndex, resp.EventTypes)
	}
}
```