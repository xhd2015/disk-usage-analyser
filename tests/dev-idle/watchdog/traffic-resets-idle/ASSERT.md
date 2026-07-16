## Expected

- Touches at T=0 and T=2s with `Timeout=3s`, then advancing to T=4s without exceeding a 3s gap, must not call `OnIdle`.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OnIdleCount != 0 {
		t.Fatalf("OnIdleCount = %d, want 0 when traffic resets idle window", resp.OnIdleCount)
	}
}
```