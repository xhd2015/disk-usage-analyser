## Expected

- Exit code 0.
- `KIND: android-avd`.
- `BREAKDOWN` includes a header row with columns `SIZE`, `NAME`, `ROLE`, `RECLAIMABLE`
  (and `NOTES` when notes are present for AVD).
- Multi-row body under BREAKDOWN; **SIZE** values are **right-aligned** (end column shared
  across rows; fixture includes both `400B` and `32B` widths).
- No `role=` prefix on human ROLE cells.
- No ANSI under default non-TTY auto color; trailing blank line; no `rm -rf`.

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
	header, data := parseBreakdownTable(resp.Stdout)
	if !strings.Contains(header, "NOTES") {
		t.Fatalf("AVD BREAKDOWN header should include NOTES when notes are present: %q", header)
	}
	if len(data) < 3 {
		t.Fatalf("BREAKDOWN expected ≥3 data rows for AVD fixture, got %d:\n%s", len(data), breakdownSection(resp.Stdout))
	}
	assertBreakdownSizeColumnAligned(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)
	assertBreakdownHasCheckboxes(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	assertNoRmRf(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
