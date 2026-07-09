## Expected

- Exit code 0.
- `KIND: android-avd`.
- Full stdout (and stderr) contain **no** ANSI escape sequences.
- BREAKDOWN still has table header + RECLAIMABLE checkboxes (`☑` for snapshots, `☐` for config).
- Trailing blank line; no `rm -rf`.

## Exit Code

- 0

```go
import (
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d (err=%v stderr=%q)", resp.ExitCode, resp.Err, resp.Stderr)
	}
	assertKindLine(t, resp.Stdout, "android-avd")
	assertNoANSI(t, resp.Stdout)
	assertNoANSI(t, resp.Stderr)
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownHasCheckboxes(t, resp.Stdout)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "snapshots", true)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "config.ini", false)
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
