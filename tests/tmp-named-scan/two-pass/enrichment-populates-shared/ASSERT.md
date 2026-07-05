## Expected

- Exactly one `named_enriched` event.
- `named_enriched.sharedSize` is 4096 (`pnpmSharedSize` + `bunSharedSize` from store+cache clones).
- `named_enriched.sharedHuman` is non-empty.
- Pass 1 `named` event has `sharedSize=0`.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedEnrichedJSON) != 1 {
		t.Fatalf("expected 1 named_enriched hit, got %d: %#v", len(resp.NamedEnrichedJSON), resp.NamedEnrichedJSON)
	}
	enriched := resp.NamedEnrichedJSON[0]
	if enriched.SharedSize != file4K {
		t.Fatalf("named_enriched sharedSize = %d, want %d (pnpm %d + bun %d)", enriched.SharedSize, file4K, enriched.PnpmSharedSize, enriched.BunSharedSize)
	}
	if enriched.SharedHuman == "" {
		t.Fatal("named_enriched sharedHuman must be non-empty")
	}
	if len(resp.NamedJSON) != 1 {
		t.Fatalf("expected 1 pass1 named hit, got %d", len(resp.NamedJSON))
	}
	if resp.NamedJSON[0].SharedSize != 0 {
		t.Fatalf("pass1 named sharedSize = %d, want 0", resp.NamedJSON[0].SharedSize)
	}
}
```