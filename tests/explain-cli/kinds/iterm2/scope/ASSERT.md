## Expected

- Exit code 0.
- Human sections in order with exact kind line: `KIND: iterm2`.
- Output includes the measured absolute iTerm2 path under the scope (and typically the scope).
- BREAKDOWN mentions iTerm2 roles/names; RECLAIMABLE ☑ for python-env / python-env-alias /
  logs, ☐ for meta / user-config.
- BREAKDOWN size DESC: iterm2env → iterm2env-3.10 → log.0.txt → version.txt → DynamicProfiles.
- `SUMMARY` / `SAFE TO RECLAIM` hardlink wording present.
- `HOW TO PURGE` is CLI-first (`disk-usage-analyser scan` and/or `du`); runnable lines use
  **`$ `**; never `rm -rf`.
- `RAW COMMANDS` includes `$ disk-usage-analyser scan` and the iTerm2 path.
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
	assertKindLine(t, resp.Stdout, "iterm2")
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("PATH/output must include measured iTerm2 %q:\n%s", req.TargetPath, resp.Stdout)
	}
	assertITerm2BreakdownMentions(t, resp.Stdout)
	assertBreakdownTableHeader(t, resp.Stdout)
	assertBreakdownNoRoleEqualsPrefix(t, resp.Stdout)
	assertBreakdownNamesInOrder(t, resp.Stdout, []string{
		"iterm2env",
		"iterm2env-3.10",
		"log.0.txt",
		"version.txt",
		"DynamicProfiles",
	})
	assertITerm2ReclaimCheckboxes(t, resp.Stdout)
	assertITerm2HardlinkWording(t, resp.Stdout)
	assertITerm2SafeToReclaim(t, resp.Stdout)
	assertITerm2CLIFirstPurge(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, req.TargetPath)
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
