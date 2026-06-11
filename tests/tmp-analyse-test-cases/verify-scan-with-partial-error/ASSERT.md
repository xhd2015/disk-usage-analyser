## Expected
- Size > 0 (accumulated values returned even with error)
- FileCount >= 1
- Progress callback was called at least once
- Error was returned (ReclaimableSize == 1)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size <= 0 {
		t.Fatalf("expected Size > 0 (accumulated before error), got %d", resp.Size)
	}
	if resp.FileCount < 1 {
		t.Fatalf("expected FileCount >= 1, got %d", resp.FileCount)
	}
	if resp.TotalSize < 1 {
		t.Fatalf("expected at least 1 progress callback, got %d", resp.TotalSize)
	}
	if resp.ReclaimableSize != 1 {
		t.Fatal("expected ScanWithProgress to return an error (partial scan)")
	}
}
```
