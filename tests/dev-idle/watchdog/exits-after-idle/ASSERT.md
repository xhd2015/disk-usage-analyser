## Expected

- After one touch and advancing the fake clock 3s with `Timeout=2s`, `OnIdle` is called exactly once.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OnIdleCount != 1 {
		t.Fatalf("OnIdleCount = %d, want 1 after idle expiry", resp.OnIdleCount)
	}
}
```