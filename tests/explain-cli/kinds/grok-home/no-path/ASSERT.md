## Expected

- Exit code 0 (PATH optional when `--kind` is set).
- `KIND: grok-home` (CLI alias is `grok`; kind id remains `grok-home`).
- `PATH:` / output reflects the measured `{HomeDir}/.grok` tree (not a missing-path error).
- BREAKDOWN covers Grok roles/names + reclaim checkboxes under that home.
- HOW TO PURGE CLI-first (scan); no `rm -rf`; `$` on runnable lines; trailing blank; no ANSI.

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
		t.Fatalf("expected exit 0 with --kind grok and no PATH (HomeDir inject), got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "grok-home")
	// Measured context is {HomeDir}/.grok
	grokDir := req.TargetPath
	if grokDir == "" || !strings.Contains(resp.Stdout, grokDir) {
		t.Fatalf("PATH/output must include measured .grok %q:\n%s", grokDir, resp.Stdout)
	}
	if req.HomeDir != "" && !strings.Contains(resp.Stdout, req.HomeDir) {
		// PATH may be absolute .grok under home; home path should still appear as prefix of .grok.
		// Soft: already required full grokDir which is under HomeDir.
	}
	assertGrokBreakdownMentions(t, resp.Stdout)
	assertGrokReclaimCheckboxes(t, resp.Stdout)
	assertGrokCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, grokDir)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
