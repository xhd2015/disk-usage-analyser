## Expected

- Exit code 0.
- `KIND: generic-dir` (not seatalk / specialized cache kinds).
- BREAKDOWN is table-shaped with header + RECLAIMABLE checkboxes.
- Children appear in **size DESC** order: `Cache` (200) → `tmp` (100) → `notes.txt` (32).
- Roles remapped: `Cache` → `cache` with `☑`; `tmp` → `tmp` with `☑`; `notes.txt` → neutral
  (`file` or similar) with `☐`.
- Bare roles (no `role=`); no ANSI; trailing blank line; no `rm -rf`.

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
	assertKindLine(t, resp.Stdout, "generic-dir")
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)

	// Size DESC by basename appearance order.
	assertBreakdownNamesInOrder(t, resp.Stdout, []string{"Cache", "tmp", "notes.txt"})

	// Role remap + reclaim checkboxes.
	cacheLine := breakdownLineForName(t, resp.Stdout, "Cache")
	tmpLine := breakdownLineForName(t, resp.Stdout, "tmp")
	notesLine := breakdownLineForName(t, resp.Stdout, "notes.txt")

	if !strings.Contains(cacheLine, "cache") {
		t.Fatalf("Cache row should use role cache (or web-cache): %q", cacheLine)
	}
	// Prefer semantic cache (not bare "directory") for reclaim signal.
	if strings.Contains(cacheLine, "directory") && !strings.Contains(cacheLine, "cache") {
		t.Fatalf("Cache row must not remain generic directory role only: %q", cacheLine)
	}
	if !strings.Contains(tmpLine, "tmp") {
		t.Fatalf("tmp row should use role tmp: %q", tmpLine)
	}
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "Cache", true)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "tmp", true)
	assertBreakdownReclaimableCheckbox(t, resp.Stdout, "notes.txt", false)

	// notes should not be reclaimable tier.
	if strings.Contains(notesLine, "☑") {
		t.Fatalf("notes.txt must not be reclaimable ☑: %q", notesLine)
	}

	assertNoANSI(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
