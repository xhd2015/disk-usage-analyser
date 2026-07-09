## Expected

- Exit code 0.
- `KIND: android-avd`.
- BREAKDOWN table header includes `RECLAIMABLE`.
- Row for `snapshots` (role `snapshot`) shows reclaimable checkbox **`☑`**.
- Rows for `userdata-qemu.img.qcow2` (`user-data`), `sdcard.img` (`sdcard`), `config.ini` (`config`)
  show **`☐`** (caution tier — not reclaimable checkbox).
- Human RECLAIMABLE column never prints the words `true` or `false`, and never ASCII `[x]` / `[ ]`.
- Bare roles (no `role=` prefix); no ANSI; trailing blank line; no `rm -rf`.

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
	assertKindLine(t, resp.Stdout, "android-avd")
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)

	// Reclaimable tier: snapshot
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "snapshots", true)
	// Caution tier: user-data, sdcard, config
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "userdata-qemu.img.qcow2", false)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "sdcard.img", false)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "config.ini", false)

	// Global: RECLAIMABLE column must not use boolean words or ASCII checkboxes.
	bd := breakdownSection(resp.Stdout)
	for _, line := range strings.Split(bd, "\n") {
		for _, f := range strings.Fields(line) {
			if f == "true" || f == "false" {
				t.Fatalf("BREAKDOWN human must not print boolean %q in RECLAIMABLE column: %q", f, line)
			}
			if f == "[x]" || f == "[ ]" {
				t.Fatalf("BREAKDOWN human must not print ASCII checkbox %q (want ☑/☐): %q", f, line)
			}
		}
	}

	assertNoANSI(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
