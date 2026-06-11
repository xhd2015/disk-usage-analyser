## Expected
- Size equals 100 (10 + 20 + 30 + 40 bytes recursively)
- FileCount equals 4 (all files in nested directories counted)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 100 {
		t.Fatalf("expected Size=100 for nested dirs, got %d", resp.Size)
	}
	if resp.FileCount != 4 {
		t.Fatalf("expected FileCount=4 for nested dirs, got %d", resp.FileCount)
	}
}
```
