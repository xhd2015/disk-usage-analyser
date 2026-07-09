## Expected

- Exit code 0.
- `KIND: homebrew-cache`.
- Human sections present; no `rm -rf`; scan in RAW COMMANDS; trailing blank line.
- HOW TO PURGE official runnable line includes `$ brew` (e.g. `$ brew cleanup`).
- RAW COMMANDS runnable lines use `$ ` prefix; no ANSI under default auto/non-TTY.

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
	assertKindLine(t, resp.Stdout, "homebrew-cache")
	assertHumanSectionsPresent(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertHowToPurgeHasDollarCommand(t, resp.Stdout, "brew")
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
