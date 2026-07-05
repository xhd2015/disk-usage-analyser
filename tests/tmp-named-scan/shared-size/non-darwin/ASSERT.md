## Expected

- Exactly one `named` hit.
- `sharedSize` is 0.
- `sharedHuman` is `"0 B"`.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.NamedJSON) != 1 {
		t.Fatalf("expected 1 named hit, got %d: %#v", len(resp.NamedJSON), resp.NamedJSON)
	}
	hit := resp.NamedJSON[0]
	if hit.SharedSize != 0 {
		t.Fatalf("sharedSize = %d, want 0 on non-darwin", hit.SharedSize)
	}
	if hit.SharedHuman != "0 B" {
		t.Fatalf("sharedHuman = %q, want %q", hit.SharedHuman, "0 B")
	}
}
```