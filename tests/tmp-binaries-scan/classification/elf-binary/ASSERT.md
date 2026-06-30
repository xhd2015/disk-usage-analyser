## Expected

- One binary hit with `kind=elf`.
- `repoName` is `elf-app`.

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
	if resp.Binaries[0].Kind != "elf" {
		t.Fatalf("kind = %q, want elf", resp.Binaries[0].Kind)
	}
	if resp.Binaries[0].RepoName != "elf-app" {
		t.Fatalf("repoName = %q, want elf-app", resp.Binaries[0].RepoName)
	}
}
```
