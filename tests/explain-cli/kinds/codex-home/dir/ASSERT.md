## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `LOGS DB`, `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: codex-home`.
- `PATH:` includes the absolute `.codex` directory.
- BREAKDOWN/summary mentions Codex roles or basenames (app-logs-db / logs_*.sqlite,
  session-logs / sessions, cache, tmp / .tmp, config / auth).
- RECLAIMABLE: logs db / sessions / cache / tmp → `☑`; config/auth → `☐`.
- **LOGS DB** section present: PATH, SIZE, **ROWS: 5**, SAMPLE (last 3, newest first) with
  newest row cues (row5 / epsilon / DEBUG).
- `SAFE TO RECLAIM`: logs/sessions/cache/tmp reclaimable language; auth/config not usually-safe.
- `HOW TO PURGE`:
  - Inspect with `$ disk-usage-analyser scan` on `.codex` (CLI-first).
  - Reclaim guidance for logs / sessions / cache / tmp; never `rm -rf`.
  - **Safe `logs_2.sqlite` reclaim** (`assertCodexHowToPurgeLogs`): quit Codex first;
    `mv`/backup of `logs_2.sqlite` + `-wal`/`-shm` **and** `sqlite3 … "DELETE FROM logs; VACUUM;"`;
    notes: diagnostic-only; not `state_5`/auth/config; Codex recreates DB; may regrow (TRACE);
    no live truncate while running.
  - Does not mark auth/config as usually-safe Removes.
  - Runnable official lines use **`$ `** prefix.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the `.codex` path.
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
	assertKindLine(t, resp.Stdout, "codex-home")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include target %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Must not fall through to generic-dir when signatures are present.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-dir" {
			t.Fatalf(".codex with Codex signatures must prefer codex-home, got generic-dir:\n%s", resp.Stdout)
		}
	}
	assertCodexBreakdownMentions(t, resp.Stdout)
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)
	assertCodexReclaimCheckboxes(t, resp.Stdout)
	assertCodexLogsDBHuman(t, resp.Stdout, codexHomeFixtureLogRows)
	assertCodexSafeToReclaim(t, resp.Stdout)
	assertCodexCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
