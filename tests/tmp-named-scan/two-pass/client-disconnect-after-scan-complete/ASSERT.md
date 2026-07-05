## Expected

- Harness reads `scan_complete` then closes the SSE body.
- Server logs must **not** contain `Error sending SSE event named_enriched` or
  `broken pipe` while workers drain after disconnect.
- Disconnect must not hang the harness.

## Errors

- Harness must not return error.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !resp.DisconnectSawScanComplete {
		t.Fatal("expected to observe scan_complete before disconnect")
	}
	if strings.Contains(resp.ServerLog, "Error sending SSE event named_enriched") {
		t.Fatalf("server must not log named_enriched SSE write errors after client disconnect:\n%s", resp.ServerLog)
	}
	if strings.Contains(resp.ServerLog, "broken pipe") {
		t.Fatalf("server must not log broken pipe after client disconnect:\n%s", resp.ServerLog)
	}
}
```