## Expected

- One binary hit with `kind=macho`.
- `repoName` is `macho-app`.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Binaries) != 1 {
		t.Fatalf("expected 1 hit, got %#v", resp.Binaries)
	}
	if resp.Binaries[0].Kind != "macho" {
		t.Fatalf("kind = %q, want macho", resp.Binaries[0].Kind)
	}
	if resp.Binaries[0].RepoName != "macho-app" {
		t.Fatalf("repoName = %q, want macho-app", resp.Binaries[0].RepoName)
	}
}
```
