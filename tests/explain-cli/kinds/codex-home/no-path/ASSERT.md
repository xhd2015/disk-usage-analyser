## Expected

- Exit code 0 (PATH optional when `--kind` is set).
- `KIND: codex-home` (CLI alias is `codex`; kind id remains `codex-home`).
- `PATH:` / output reflects the measured `{HomeDir}/.codex` tree (not a missing-path error).
- BREAKDOWN covers Codex roles/names + reclaim checkboxes under that home.
- LOGS DB present with ROWS: 5.
- HOW TO PURGE CLI-first (scan) + safe `logs_2.sqlite` reclaim (quit/mv/sqlite3); no `rm -rf`;
  `$` on runnable lines; trailing blank; no ANSI.

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
		t.Fatalf("expected exit 0 with --kind codex and no PATH (HomeDir inject), got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "codex-home")
	// Measured context is {HomeDir}/.codex
	codexDir := req.TargetPath
	if codexDir == "" || !strings.Contains(resp.Stdout, codexDir) {
		t.Fatalf("PATH/output must include measured .codex %q:\n%s", codexDir, resp.Stdout)
	}
	assertCodexBreakdownMentions(t, resp.Stdout)
	assertCodexReclaimCheckboxes(t, resp.Stdout)
	assertCodexLogsDBHuman(t, resp.Stdout, codexHomeFixtureLogRows)
	assertCodexCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, codexDir)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
