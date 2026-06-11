## Expected
- p1 (rebootSafe, cur=200, acc=1000): totalSize=1200, reclaimableSize=1200
- p2 (not rebootSafe, cur=500, acc=1200/1200): totalSize=1700, reclaimableSize=1200

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 1200 {
		t.Fatalf("expected totalSize=1200 for p1 (1000+200), got %d", resp.Size)
	}
	if resp.FileCount != 1200 {
		t.Fatalf("expected reclaimableSize=1200 for p1 (1000+200, rebootSafe), got %d", resp.FileCount)
	}
	if resp.TotalSize != 1700 {
		t.Fatalf("expected totalSize=1700 for p2 (1200+500), got %d", resp.TotalSize)
	}
	if resp.ReclaimableSize != 1200 {
		t.Fatalf("expected reclaimableSize=1200 for p2 (1200+0, not rebootSafe), got %d", resp.ReclaimableSize)
	}
}
```
