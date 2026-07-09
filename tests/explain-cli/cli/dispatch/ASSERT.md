## Expected

- No harness error from `run.RunWithOptions`.
- Web server hook is not invoked.
- Stdout is a successful human explain for the fixture path (includes `PATH:` and `KIND:`).
- Stdout ends with a trailing blank line.

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
	if resp.ServerWasStarted {
		t.Fatal("web server must not start for explain dispatch")
	}
	if !strings.Contains(resp.Stdout, "PATH:") {
		t.Fatalf("stdout missing PATH: header:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, "KIND:") {
		t.Fatalf("stdout missing KIND: header:\n%s", resp.Stdout)
	}
	if !strings.Contains(resp.Stdout, req.FixtureDir) {
		t.Fatalf("stdout missing fixture path %q:\n%s", req.FixtureDir, resp.Stdout)
	}
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
