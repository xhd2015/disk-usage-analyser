## Expected
- TotalSize = 1000 + 2048 + 500 = 3548 (swap included)
- ReclaimableSize = 1000 only (swap excluded even though rebootSafe=true)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalSize != 3548 {
		t.Fatalf("expected TotalSize=3548 (includes swap), got %d", resp.TotalSize)
	}
	if resp.ReclaimableSize != 1000 {
		t.Fatalf("expected ReclaimableSize=1000 (swap excluded), got %d", resp.ReclaimableSize)
	}
}
```
