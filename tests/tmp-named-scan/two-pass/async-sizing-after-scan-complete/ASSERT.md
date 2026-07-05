## Expected

- External lifecycle probe exits 0: `tmpfiles.Scan` does not return until async
  `emitNamedDirSized` goroutines have called `OnNamedHit`.
- No post-scan callbacks (`postScan=0`, `atReturn == total`).

## Errors

- Probe or harness must surface lifecycle violation with probe output.

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Helper()
	if err != nil {
		if strings.Contains(err.Error(), "async sizing finished after Scan returned") {
			t.Fatalf("lifecycle bug reproduced: %v", err)
		}
		t.Fatalf("unexpected error: %v", err)
	}
}
```