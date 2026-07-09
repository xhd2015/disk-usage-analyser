## Expected

- Root help is printed (no web server start).
- Subcommand list includes `scan` and does **not** list a standalone `inspect` subcommand.
- Mentions of inspect, if any, must only refer to `scan --inspect` (not `inspect [FILE]` as a peer subcommand).

## Exit Code

- 0 (nil error) or help-specific non-fatal return

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ServerWasStarted {
		t.Fatal("root -h must not start web server")
	}
	// Prefer stdout; some flag libraries write help to stderr.
	out := resp.Stdout + resp.Stderr
	if out == "" && resp.Err != nil {
		// help may surface as returned error containing usage text
		out = resp.Err.Error()
	}
	if !strings.Contains(out, "scan") {
		t.Fatalf("root help should mention scan:\n%s", out)
	}
	// Reject standalone inspect subcommand lines like "inspect [FILE]" as a peer of scan.
	for _, line := range strings.Split(out, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "inspect ") || trim == "inspect" {
			t.Fatalf("root help must not list inspect as a subcommand; line %q\n%s", line, out)
		}
	}
}
```
