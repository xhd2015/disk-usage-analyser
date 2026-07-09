## Expected

- Exit code 0.
- Kind is `grok-home` (`.grok` ContextRoot preferred over `generic-file`).
- Output includes the explained file path (`config.toml`) and/or the `.grok` context path.
- Human sections present; no `rm -rf`; trailing blank line.
- HOW TO PURGE still grok-oriented (scan / reclaim guidance); RAW COMMANDS use `$ ` prefix;
  no ANSI under default auto.

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
	assertKindLine(t, resp.Stdout, "grok-home")
	assertHumanSectionsPresent(t, resp.Stdout)
	if !strings.Contains(resp.Stdout, req.TargetPath) {
		t.Fatalf("output must include target file path %q:\n%s", req.TargetPath, resp.Stdout)
	}
	// ContextRoot should be the .grok directory (parent of config.toml).
	grokDir := filepath.Dir(req.TargetPath)
	if !strings.Contains(resp.Stdout, grokDir) {
		t.Fatalf("output must include .grok context root %q:\n%s", grokDir, resp.Stdout)
	}
	// Must not mis-classify as generic-file when inside .grok with signatures.
	for _, line := range strings.Split(resp.Stdout, "\n") {
		if strings.TrimSpace(line) == "KIND: generic-file" {
			t.Fatalf("file under .grok must prefer grok-home, got generic-file:\n%s", resp.Stdout)
		}
	}
	how := sectionBody(resp.Stdout, "HOW TO PURGE")
	if !strings.Contains(strings.ToLower(how), "disk-usage-analyser scan") {
		t.Fatalf("file-under-.grok explain should still include scan HOW TO PURGE:\n%s", how)
	}
	assertHowToPurgeOfficialDollarPrefix(t, resp.Stdout)
	assertRawCommandsDollarPrefix(t, resp.Stdout)
	assertContainsScanCommand(t, resp.Stdout, "")
	assertNoRmRf(t, resp.Stdout)
	assertNoANSI(t, resp.Stdout)
	stdoutEndsWithBlankLine(t, resp.Stdout)
}
```
