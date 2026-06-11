## Expected
- Size equals 600 (sum of 100 + 200 + 300)
- FileCount equals 3 (three files total, directories not counted)

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Size != 600 {
		t.Fatalf("expected Size=600, got %d", resp.Size)
	}
	if resp.FileCount != 3 {
		t.Fatalf("expected FileCount=3, got %d", resp.FileCount)
	}
}
```
