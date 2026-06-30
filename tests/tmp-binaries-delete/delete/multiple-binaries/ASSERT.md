## Expected

- HTTP status 200.
- `deleted` contains two paths.
- `failed` is empty.
- `freedSize` equals sum of both file sizes.
- Neither file exists after delete.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, body=%s", resp.StatusCode, resp.Body)
	}
	if len(resp.Result.Deleted) != 2 {
		t.Fatalf("deleted = %#v, want 2", resp.Result.Deleted)
	}
	if len(resp.Result.Failed) != 0 {
		t.Fatalf("failed = %#v", resp.Result.Failed)
	}
	if resp.Result.FreedSize <= 0 {
		t.Fatalf("freedSize = %d", resp.Result.FreedSize)
	}
	for _, p := range resp.Result.Deleted {
		if resp.FileExists[p] {
			t.Fatalf("file still exists: %s", p)
		}
	}
}
```
