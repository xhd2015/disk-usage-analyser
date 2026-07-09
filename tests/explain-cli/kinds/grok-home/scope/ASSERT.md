## Expected

- Exit code 0.
- Human sections in order with exact kind line: `KIND: grok-home`.
- Output includes the measured absolute `.grok` path under the scope (and typically the scope).
- BREAKDOWN mentions Grok roles/names; RECLAIMABLE ☑ for downloads/sessions/marketplace/logs,
  ☐ for config/auth.
- BREAKDOWN size DESC: downloads → sessions → marketplace-cache → logs.
- `SAFE TO RECLAIM` present with reclaim language for installer/cache; not auth/config.
- `HOW TO PURGE` is CLI-first (`disk-usage-analyser scan`); runnable lines use **`$ `**; never `rm -rf`.
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
		t.Fatalf("PATH/output must include measured .grok %q:\n%s", req.TargetPath, resp.Stdout)
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
