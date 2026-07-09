## Expected

- Exit code 0.
- `KIND: generic-qcow2` (not `android-avd` — no parent AVD context).
- Human sections present; no `rm -rf`; trailing blank line.
- RAW COMMANDS use `$ ` prefix; official HOW TO PURGE lines that are runnable use `$ `;
  comment-only official text is allowed without `$`.
- No ANSI under default auto/non-TTY.

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
	assertKindLine(t, resp.Stdout, "generic-qcow2")
	assertHumanSectionsPresent(t, resp.Stdout)
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: android-avd" {
			t.Fatalf("lone qcow2 without AVD context must not be android-avd:\n%s", resp.Stdout)
		}
	}
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
