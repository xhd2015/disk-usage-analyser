## Expected

- Exit code 0.
- Kind is `codex-home` (`.codex` ContextRoot preferred over `generic-file`).
- Output includes the explained file path (`config.toml`) and/or the `.codex` context path.
- Human sections present; LOGS DB still shown when fixture has logs sqlite.
- HOW TO PURGE still codex-oriented (scan / reclaim guidance); RAW COMMANDS use `$ ` prefix;
  no `rm -rf`; trailing blank; no ANSI under default auto.

## Exit Code

- 0

```go
import (
	"path/filepath"
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
	assertKindLine(t, resp.Stdout, "codex-home")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output must include target file path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// ContextRoot should be the .codex directory (parent of config.toml).
	codexDir := filepath.Dir(req.TargetPath)
	if !strings.Contains(resp.Stdout, codexDir) {
		t.Fatalf("output must include .codex context root %q:\n%s", codexDir, resp.Stdout)
	}
	// Must not mis-classify as generic-file when inside .codex with signatures.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-file" {
			t.Fatalf("file under .codex must prefer codex-home, got generic-file:\n%s", resp.Stdout)
		}
	}
	assertCodexLogsDBHuman(t, resp.Stdout, codexHomeFixtureLogRows)
	how := sectionBody(resp.Stdout, "HOW TO PURGE")
	if !strings.Contains(strings.ToLower(how), "disk-usage-analyser scan") {
		t.Fatalf("file-under-.codex explain should still include scan HOW TO PURGE:\n%s", how)
	}
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
