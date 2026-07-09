## Expected

- No harness error from `run.RunWithOptions`.
- Web server hook is not invoked.
- Stdout is inspect tree view: `PATH:`, `MIN:`, `SOURCE:`, depth-1 children, trailing blank line.
- No TOP section (no query flags).

## Exit Code

- 0 (via nil error from run)

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("unexpected run error: %v", resp.Err)
	}
	if resp.ServerWasStarted {
		t.Fatal("web server must not start for scan --inspect dispatch")
	}
	out := resp.Stdout
	for _, want := range []string{"PATH:", "MIN:", "MAX-DEPTH:", "SOURCE:", "huge.bin", "mid.bin", "big/"} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "deep.bin") {
		t.Fatalf("default inspect max-depth 1 must hide deep.bin:\n%s", out)
	}
	if strings.Contains(out, "THRESHOLD:") {
		t.Fatalf("must use MIN: not THRESHOLD:\n%s", out)
	}
	stdoutHasNoTopSection(t, out)
	stdoutEndsWithBlankLine(t, out)
}
```
