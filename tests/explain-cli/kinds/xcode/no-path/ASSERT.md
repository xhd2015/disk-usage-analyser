## Expected

- Exit code 0 (PATH optional when `--kind` is set).
- `KIND: xcode`.
- `PATH:` reflects the injected home/scope (`req.HomeDir` / `req.TargetPath`), not a missing-path error.
- BREAKDOWN covers the Xcode pack under that scope (roles/names + reclaim checkboxes).
- HOW TO PURGE CLI-first; no `rm -rf`; `$` on runnable lines; trailing blank; no ANSI.

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
		t.Fatalf("expected exit 0 with --kind and no PATH (HomeDir inject), got %d (err=%v stderr=%q stdout=%q)",
			resp.ExitCode, resp.Err, resp.Stderr, resp.Stdout)
	}
	assertHumanSectionsPresent(t, resp.Stdout)
	assertHumanSectionOrder(t, resp.Stdout)
	assertKindLine(t, resp.Stdout, "xcode")
	scope := req.HomeDir
	if scope == "" {
		scope = req.TargetPath
	}
	if scope == "" || !strings.Contains(resp.Stdout, scope) {
		t.Fatalf("PATH/output must include injected home/scope %q:\n%s", scope, resp.Stdout)
	}
	assertXcodeBreakdownMentions(t, resp.Stdout)
	assertXcodeReclaimCheckboxes(t, resp.Stdout)
	assertXcodeCLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, scope)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
