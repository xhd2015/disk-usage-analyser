## Expected

- Exit code 0 (PATH optional when `--kind` is set).
- `KIND: iterm2`.
- `PATH:` / output reflects the measured `{HomeDir}/Library/Application Support/iTerm2`
  tree (not a missing-path error).
- BREAKDOWN covers iTerm2 roles/names + reclaim checkboxes under that home.
- HOW TO PURGE CLI-first (`scan` / `du`); no `rm -rf`; `$` on runnable lines; trailing blank;
  no ANSI.

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
		t.Fatalf("expected exit 0 with --kind iterm2 and no PATH (HomeDir inject), got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "iterm2")
	// Measured context is {HomeDir}/Library/Application Support/iTerm2
	iterm2Dir := req.TargetPath
	if iterm2Dir == "" || !strings.Contains(resp.Stdout, iterm2Dir) {
		t.Fatalf("PATH/output must include measured iTerm2 %q:\n%s", iterm2Dir, resp.Stdout)
	}
	if req.HomeDir != "" && !strings.Contains(iterm2Dir, req.HomeDir) {
		t.Fatalf("measured iTerm2 %q should be under HomeDir %q", iterm2Dir, req.HomeDir)
	}
	assertITerm2BreakdownMentions(t, resp.Stdout)
	assertITerm2ReclaimCheckboxes(t, resp.Stdout)
	assertITerm2CLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, iterm2Dir)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
