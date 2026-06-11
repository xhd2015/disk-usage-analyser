## Expected
- ScanWithProgress returns size=1000, fileCount=4
- Progress callback was called at least 4 times (once per file)
- Final progress size equals 1000

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 1000 {
		t.Fatalf("expected Size=1000, got %d", resp.Size)
	}
	if resp.FileCount != 4 {
		t.Fatalf("expected FileCount=4, got %d", resp.FileCount)
	}
	if resp.TotalSize < 4 {
		t.Fatalf("expected at least 4 progress callbacks, got %d", resp.TotalSize)
	}
	if resp.ReclaimableSize != 1000 {
		t.Fatalf("expected final progress size=1000, got %d", resp.ReclaimableSize)
	}
}
```
