## Expected

- At least one `named` event.
- Every pass 1 `named` event has `sharedSize=0` and `sharedHuman="0 B"`.
- If `named_enriched` events exist, every `named` appears before the first `named_enriched` in event order.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedJSON) == 0 {
		t.Fatal("expected at least one named SSE event")
	}
	for i, hit := range resp.NamedJSON {
		if hit.SharedSize != 0 {
			t.Fatalf("named event %d sharedSize = %d, want 0 in pass 1", i, hit.SharedSize)
		}
		if hit.SharedHuman != "0 B" {
			t.Fatalf("named event %d sharedHuman = %q, want %q in pass 1", i, hit.SharedHuman, "0 B")
		}
	}
	if resp.HasNamedEnriched && !eventLastBefore(resp.EventTypes, "named", "named_enriched") {
		t.Fatalf("expected all named events before first named_enriched, got events: %v", resp.EventTypes)
	}
}
```