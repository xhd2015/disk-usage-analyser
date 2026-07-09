## Expected

- Exit code 0.
- `KIND: android-avd`.
- Stdout contains ANSI (color forced).
- In **BREAKDOWN**:
  - reclaimable role `snapshot` is **green** (`\x1b[32m` or bold green).
  - caution roles `user-data`, `sdcard`, `config` are **yellow** (`\x1b[33m` or bold yellow).
  - reclaimable checkbox **`☑` is green-wrapped** (same SGR family as reclaimable ROLE).
  - non-reclaimable **`☐` is not** green/yellow-wrapped.
- Green base commands on HOW TO PURGE / RAW (`emulator`, `disk-usage-analyser`, …) remain allowed.
- Checkboxes still present (`☑`/`☐`); trailing blank line; no `rm -rf`.

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
	if !containsANSI(resp.Stdout) {
		t.Fatalf("expected ANSI with --color=always:\n%s", resp.Stdout)
	}

	// ROLE cell colors (BREAKDOWN section only).
	assertBreakdownRoleGreen(t, resp.Stdout, "snapshot")
	assertBreakdownRoleYellow(t, resp.Stdout, "user-data")
	assertBreakdownRoleYellow(t, resp.Stdout, "sdcard")
	assertBreakdownRoleYellow(t, resp.Stdout, "config")

	// Reclaimable ☑ is green; non-reclaimable ☐ stays plain.
	assertBreakdownReclaimableCheckboxGreen(t, resp.Stdout)
	assertBreakdownNonReclaimableCheckboxNotColored(t, resp.Stdout)

	// Checkboxes still present (visible after strip).
	assertBreakdownHasCheckboxes(t, resp.Stdout)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "snapshots", true)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "config.ini", false)

	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
