## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: grok-home`.
- `PATH:` includes the absolute `.grok` directory.
- BREAKDOWN/summary mentions Grok roles or basenames (installer-cache / downloads,
  session-logs / sessions, cache / marketplace-cache, logs, config / auth).
- BREAKDOWN size DESC roughly: downloads → sessions → marketplace-cache → logs → config/auth.
- RECLAIMABLE: downloads / sessions / marketplace / logs → `☑`; config/auth → `☐`.
- `SAFE TO RECLAIM`: installer/cache/logs reclaimable language; auth/config not usually-safe purge.
- `HOW TO PURGE`:
  - Inspect with `$ disk-usage-analyser scan` on `.grok` / downloads / sessions (CLI-first).
  - Reclaim guidance for downloads / caches; never `rm -rf`.
  - Does not mark auth/config as usually-safe Removes.
  - Runnable official lines use **`$ `** prefix.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the `.grok` path.
- Full stdout: no `rm -rf`; trailing blank line; no ANSI under default non-TTY auto color.

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
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "grok-home")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include target %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Must not fall through to generic-dir when signatures are present.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-dir" {
			t.Fatalf(".grok with Grok signatures must prefer grok-home, got generic-dir:\n%s", resp.Stdout)
		}
	}
	assertGrokBreakdownMentions(t, resp.Stdout)
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)
	assertBreakdownNamesInOrder(t, resp.Stdout, []string{
		"downloads",
		"sessions",
		"marketplace-cache",
		"logs",
	})
	assertGrokReclaimCheckboxes(t, resp.Stdout)
	assertGrokSafeToReclaim(t, resp.Stdout)
	assertGrokCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
