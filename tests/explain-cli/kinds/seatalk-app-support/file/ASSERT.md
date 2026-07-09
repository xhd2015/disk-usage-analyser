## Expected

- Exit code 0.
- Kind is `seatalk-app-support` (SeaTalk ContextRoot preferred over `generic-file`).
- Output includes the explained file path (`main_1.sqlite`).
- Human sections present; no `rm -rf`; trailing blank line.
- HOW TO PURGE still seatalk-oriented (osascript quit); RAW COMMANDS use `$ ` prefix;
  no ANSI under default auto.

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
	assertKindLine(t, resp.Stdout, "seatalk-app-support")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output must include target file path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Must not mis-classify as generic-file when inside SeaTalk Application Support.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-file" {
			t.Fatalf("file under SeaTalk must prefer seatalk-app-support, got generic-file:\n%s", resp.Stdout)
		}
	}
	how := sectionBody(resp.Stdout, "HOW TO PURGE")
	if !strings.Contains(how, "osascript") {
		t.Fatalf("file-under-SeaTalk explain should still include osascript HOW TO PURGE:\n%s", how)
	}
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
