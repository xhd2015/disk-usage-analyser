## Expected
- Size equals 0
- FileCount equals 0
- No error returned

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 0 {
		t.Fatalf("expected Size=0 for empty dir, got %d", resp.Size)
	}
	if resp.FileCount != 0 {
		t.Fatalf("expected FileCount=0 for empty dir, got %d", resp.FileCount)
	}
}
```
