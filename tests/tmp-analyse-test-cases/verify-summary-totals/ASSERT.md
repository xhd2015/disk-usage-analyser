## Expected
- TotalSize equals 3800 (sum of all location sizes)
- ReclaimableSize equals 1800 (sum of only rebootSafe=true locations: 1000+500+300)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.TotalSize != 3800 {
		t.Fatalf("expected TotalSize=3800, got %d", resp.TotalSize)
	}
	if resp.ReclaimableSize != 1800 {
		t.Fatalf("expected ReclaimableSize=1800 (only rebootSafe), got %d", resp.ReclaimableSize)
	}
}
```
