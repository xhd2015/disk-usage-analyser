## Expected

- Exit code 0.
- Successful human explain (includes `SAFE TO RECLAIM`).
- Entire stdout must not contain `rm -rf` (case-insensitive).
- Prefer safer guidance language (e.g. delete/remove snapshots via tooling) without shell nukes.

## Exit Code

- 0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d (err=%v stderr=%q)", resp.ExitCode, resp.Err, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, "SAFE TO RECLAIM") {
		t.Fatalf("expected SAFE TO RECLAIM section so safety check is meaningful:\n%s", resp.Stdout)
	}
	assertNoRmRf(t, resp.Stdout)
	// Also reject close variants that tests treat as equivalent destructive one-liners.
	lower := strings.ToLower(resp.Stdout)
	for _, bad := range []string{"rm -rf", "rm -fr", "rm -r -f", "rm -f -r"} {
		if strings.Contains(lower, bad) {
			t.Fatalf("stdout must not contain destructive pattern %q:\n%s", bad, resp.Stdout)
		}
	}
	// Default human path is non-TTY auto: must not inject ANSI into safety-critical text.
	assertNoANSI(t, resp.Stdout)
}
```
