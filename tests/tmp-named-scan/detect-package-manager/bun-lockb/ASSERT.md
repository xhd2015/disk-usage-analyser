## Expected

- Exactly one `named` hit.
- `packageManager` is `bun`.

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
	if resp.NamedJSON[0].PackageManager != "bun" {
		t.Fatalf("packageManager = %q, want bun", resp.NamedJSON[0].PackageManager)
	}
}
```