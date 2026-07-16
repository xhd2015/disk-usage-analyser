## Expected

- `OnIdle` is never called when `Timeout` is zero, even after touch and long clock advance.

## Side Effects

- None beyond watch lifecycle (Start/Stop test hooks).

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.OnIdleCount != 0 {
		t.Fatalf("OnIdleCount = %d, want 0 when Timeout=0", resp.OnIdleCount)
	}
}
```