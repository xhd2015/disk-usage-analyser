## Expected

- Exactly one `named` hit.
- SSE JSON includes `hasPackageJson: false`.
- `packageManager` is `unknown`.

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
	objs := namedEventDataObjects(t, resp.SSEOutput)
	if len(objs) != 1 {
		t.Fatalf("expected 1 named event object, got %d", len(objs))
	}
	hasPkg, ok := objs[0]["hasPackageJson"].(bool)
	if !ok {
		t.Fatalf("hasPackageJson key missing from named event: %#v", objs[0])
	}
	if hasPkg {
		t.Fatalf("hasPackageJson = true, want false")
	}
	if resp.NamedJSON[0].PackageManager != "unknown" {
		t.Fatalf("packageManager = %q, want unknown", resp.NamedJSON[0].PackageManager)
	}
}
```