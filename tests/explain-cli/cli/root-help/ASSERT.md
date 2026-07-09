## Expected

- Root help is printed (no web server start).
- Subcommand list includes `explain` and documents a PATH argument (e.g. `explain [PATH]`).

## Exit Code

- 0 (nil error from run)

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
	out := resp.Stdout + resp.Stderr
	if out == "" && resp.Err != nil {
		out = resp.Err.Error()
	}
	if !strings.Contains(out, "explain") {
		t.Fatalf("root help should mention explain:\n%s", out)
	}
	// Prefer the documented form "explain [PATH]" when present; also accept a line that
	// starts with explain and mentions PATH nearby.
	if !strings.Contains(out, "explain [PATH]") && !strings.Contains(out, "explain PATH") {
		// Fall back: a subcommand line containing both tokens
		found := false
		for _, line := range strings.Split(out, "\n") {
			trim := strings.TrimSpace(line)
			if strings.Contains(trim, "explain") && strings.Contains(strings.ToUpper(trim), "PATH") {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("root help should list explain with PATH (e.g. 'explain [PATH]'):\n%s", out)
		}
	}
}
```
