## Expected

- Exit code 0.
- Human sections in order: `PATH:`, `KIND:`, `TOTAL:`, `CONFIDENCE:`, `SUMMARY`, `BREAKDOWN`,
  `SAFE TO RECLAIM`, `HOW TO PURGE`, `RAW COMMANDS`.
- Exact kind line: `KIND: iterm2`.
- `PATH:` includes the absolute iTerm2 Application Support directory.
- BREAKDOWN/summary mentions iTerm2 roles or basenames (`python-env` / `iterm2env`,
  `python-env-alias` / `iterm2env-3.10`, `logs` / `log.0.txt`, `meta` / `version.txt`,
  `user-config` / `DynamicProfiles`).
- BREAKDOWN size DESC roughly: `iterm2env` → `iterm2env-3.10` → `log.0.txt` →
  `version.txt` → `DynamicProfiles`.
- RECLAIMABLE: `python-env` / `python-env-alias` / `logs` → `☑`; `meta` / `user-config` → `☐`.
- `SUMMARY` or `SAFE TO RECLAIM` mentions hardlink / shared blocks among `iterm2env*` and
  that logical TOTAL can overstate freeable space; confirm with `du` on parent.
- `SAFE TO RECLAIM`: python env / logs reclaimable language; user-config / DynamicProfiles
  not usually-safe-only purge.
- `HOW TO PURGE`:
  - CLI-first: `$ disk-usage-analyser scan` on iTerm2 and/or `$ du -sh` parent / envs.
  - Optional logs / all `iterm2env*` after quit guidance.
  - Never `rm -rf`.
  - Does not mark DynamicProfiles / Scripts as usually-safe Removes.
  - Runnable official lines use **`$ `** prefix.
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
		t.Fatalf("PATH/output must include target %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// Must not fall through to generic-dir when signatures are present.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-dir" {
			t.Fatalf("iTerm2 App Support tree must prefer iterm2, got generic-dir:\n%s", resp.Stdout)
		}
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
