## Expected

- Exit code 0.
- `KIND: generic-file`.
- Human sections present; RAW COMMANDS includes `$` scan; no `rm -rf`; trailing blank line.
- HOW TO PURGE official runnable lines (if any) use `$ `; no ANSI under default auto.

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
	assertKindLine(t, resp.Stdout, "generic-file")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output missing path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	assertContainsScanCommand(t, resp.Stdout, "")
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
