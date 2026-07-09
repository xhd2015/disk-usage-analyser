## Expected

- Exit code 0.
- One JSON object (+ trailing blank line); no ANSI; no `rm -rf`.
- `kind` is `"android-avd"`.
- `breakdown` is a non-empty array, **sorted by `size` descending** (nonincreasing).
- Every breakdown entry has **`reclaimable` as a JSON boolean** (not a string / checkbox).
- Role tiers: `snapshot` → `reclaimable=true`; `user-data`, `sdcard`, `config` → `false`.
- JSON must not use human checkbox glyphs (`☑` / `☐` or legacy `[x]` / `[ ]`) as the reclaim signal.

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
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)

	payload := parseJSONObject(t, resp.Stdout)
	if jsonStringField(t, payload, "kind") != "android-avd" {
		t.Fatalf("json kind: want android-avd, got %q", jsonStringField(t, payload, "kind"))
	}

	breakdown := parseJSONBreakdown(t, payload)
	assertJSONBreakdownSortedDesc(t, breakdown)
	assertJSONBreakdownReclaimableBools(t, breakdown)

	assertJSONBreakdownRoleReclaimable(t, breakdown, "snapshot", true)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "user-data", false)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "sdcard", false)
	assertJSONBreakdownRoleReclaimable(t, breakdown, "config", false)

	// Boolean field is the JSON contract — checkbox glyphs are human-only.
	line := firstJSONObjectLine(t, resp.Stdout)
	if strings.Contains(line, `"reclaimable":"[x]"`) || strings.Contains(line, `"reclaimable":"[ ]"`) ||
		strings.Contains(line, `"reclaimable":"☑"`) || strings.Contains(line, `"reclaimable":"☐"`) {
		t.Fatalf("json reclaimable must be bool, not checkbox string: %s", line)
	}
	if strings.Contains(line, "☑") || strings.Contains(line, "☐") {
		t.Fatalf("json must not embed human RECLAIMABLE glyphs ☑/☐: %s", line)
	}
}
```
