## Expected

- Harness observes at least one `repo` SSE event before disconnect.
- `resp.DisconnectAborted` is true (client cancelled after first repo event).
- Disconnect returns without hanging the test.

## Errors

- No harness error is returned.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.DisconnectAborted {
		t.Fatal("expected disconnect to abort scan after first repo event")
	}
}
```