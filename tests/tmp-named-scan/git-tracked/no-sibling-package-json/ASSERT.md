## Expected

- Exactly one `named_enriched` hit.
- `gitTracked` is `false`.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedEnrichedJSON) != 1 {
		t.Fatalf("expected 1 named_enriched hit, got %d", len(resp.NamedEnrichedJSON))
	}
	assertNamedEnrichedGitTracked(t, resp, false)
}
```