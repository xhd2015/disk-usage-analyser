## Expected

- Exit code 0.
- Stdout contains green ANSI around base command token **`go`** (HOW TO PURGE and/or RAW).
- Stdout also greens **`disk-usage-analyser`** on the scan line when present.
- Green uses SGR green (`\x1b[32m` or bold green `\x1b[1;32m` / similar).
- Shell prompt **`$` is not green**.
- Command remainder after the base token is not required to be colored.

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
	if !containsANSI(resp.Stdout) {
		t.Fatalf("expected ANSI color sequences with --color=always, got plain output:\n%s", resp.Stdout)
	}
	// Primary: green "go" from go clean -cache.
	assertGreenBaseCommand(t, resp.Stdout, "go")
	// Scan line base command should also be green when RAW COMMANDS is shown.
	if strings.Contains(resp.Stdout, "disk-usage-analyser") {
		assertGreenBaseCommand(t, resp.Stdout, "disk-usage-analyser")
	}
	assertNoRmRf(t, resp.Stdout)
}
```
