## Expected

- Exit code 0.
- Human sections in order with exact kind line: `KIND: codex-home`.
- Output includes the measured absolute `.codex` path under the scope (and typically the scope).
- BREAKDOWN mentions Codex roles/names; RECLAIMABLE ☑ for logs/sessions/cache/tmp, ☐ for config/auth.
- **LOGS DB** present with ROWS: 5 and SAMPLE last 3.
- `SAFE TO RECLAIM` present with reclaim language for logs/cache; not auth/config.
- `HOW TO PURGE` is CLI-first (`disk-usage-analyser scan`); includes safe `logs_2.sqlite`
  reclaim (quit Codex; mv backup + wal/shm and/or sqlite3 DELETE/VACUUM; not state_5/auth);
  runnable lines use **`$ `**; never `rm -rf`.
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
		t.Fatalf("PATH/output must include measured .codex %q:\n%s", req.TargetPath, resp.Stdout)
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
